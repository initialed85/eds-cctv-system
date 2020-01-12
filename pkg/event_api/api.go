package event_api

import (
	"encoding/json"
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/api"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"log"
	"net/http"
	"time"
)

type errorResponse struct {
	Error string `json:"error"`
}

type API struct {
	readStore, writeStore event_store.Store
	api                   *api.API
}

func New(path string, port int) *API {
	a := API{
		readStore:  event_store.NewStore(path),
		writeStore: event_store.NewStore(path),
		api:        api.New(port),
	}

	a.api.AddHandler("/events", a.handleEvents)
	a.api.AddHandler("/events_by_date", a.handleEventsByDate)

	a.updateStore()

	return &a
}

func (a *API) updateStore() {
	err := a.readStore.Read()
	if err != nil {
		log.Printf("failed to update readStore because: %v", err)
	}
}

func (a *API) handleResponse(v interface{}, w http.ResponseWriter) {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		message := fmt.Sprintf("failed to marshal %+v because %v", v, err)

		b, err = json.MarshalIndent(errorResponse{message}, "", "    ")

		log.Printf(message)

		return
	}

	_, err = fmt.Fprintf(w, "%v\n", string(b))
	if err != nil {
		log.Printf("failed to write response because: %v", err)

		return
	}
}

func (a *API) handleEvents(w http.ResponseWriter, r *http.Request) {
	a.updateStore()

	a.handleResponse(a.readStore.GetAll(), w)
}

func (a *API) handleEventsByDate(w http.ResponseWriter, r *http.Request) {
	a.updateStore()

	a.handleResponse(a.readStore.GetAllByDate(), w)
}

func (a *API) AddEvent(timestamp time.Time, cameraName, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath string) error {
	event := event_store.NewEvent(timestamp, cameraName, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath)

	a.writeStore.Add(event)

	err := a.writeStore.Append()
	if err != nil {
		return fmt.Errorf("failed to write %+v because %v", event, err)
	}

	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal %+v because %v", event, err)
	}

	a.api.AddMessage(string(b))

	return nil
}

func (a *API) Start() {
	go func() {
		err := a.api.Listen()

		log.Fatalf("failed to listen because: %v", err)
	}()
}

func (a *API) Stop() {
	_ = a.api.Stop()
}
