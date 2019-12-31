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
	store event_store.Store
	api   *api.API
}

func New(path string, port int) *API {
	a := API{
		store: event_store.NewStore(path),
		api:   api.New(port),
	}

	a.api.AddHandler("/events", a.handleEvents)
	a.api.AddHandler("/events_by_date", a.handleEventsByDate)

	return &a
}

func (a *API) updateStore() {
	err := a.store.Read()
	if err != nil {
		log.Printf("failed to update store because: %v", err)
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

	a.handleResponse(a.store.GetAll(), w)
}

func (a *API) handleEventsByDate(w http.ResponseWriter, r *http.Request) {
	a.updateStore()

	a.handleResponse(a.store.GetAllByDate(), w)
}

func (a *API) AddEvent(timestamp time.Time, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath string) error {
	event := event_store.NewEvent(timestamp, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath)

	a.store.Add(event)

	err := a.store.Write()
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
