package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Reason     string `json:"reason"`
}

func checkOrigin(_ *http.Request) bool {
	return true // TODO: this is probably a bad idea
}

func writeError(statusCode int, w http.ResponseWriter, reason string) {
	buf, err := json.MarshalIndent(ErrorResponse{statusCode, reason}, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err = fmt.Fprintf(w, "%v\n", string(buf))
	if err != nil {
		log.Fatal(err)
	}
}

type API struct {
	mux      *http.ServeMux
	server   http.Server
	ug       websocket.Upgrader
	mu       sync.Mutex
	messages []string
	shutdown bool
}

func New(port int) *API {
	mux := http.NewServeMux()

	a := &API{
		mux: mux,
		server: http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: mux,
		},
		ug:       websocket.Upgrader{CheckOrigin: checkOrigin},
		messages: make([]string, 0),
		shutdown: false,
	}

	a.mux.HandleFunc("/stream", a.handleStream)

	return a
}

func (a *API) handleStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleStream called")

	c, err := a.ug.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade because: %v", err)

		return
	}

	log.Printf("upgraded to WebSocket")

	defer func() {
		err = c.Close()
		if err != nil {
		}
	}()

	for {
		if len(a.messages) == 0 {
			time.Sleep(time.Millisecond)

			continue
		}

		log.Printf("sending %v messages", len(a.messages))

		err = nil

		a.mu.Lock()

		for _, message := range a.messages {
			err = c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("bailing on message loop because: %v", err)

				break
			}
		}

		a.messages = nil

		a.mu.Unlock()

		if err != nil {
			log.Printf("bailing on handleStream because: %v", err)

			break
		}
	}
}

func (a *API) AddMessage(message string) {
	a.mu.Lock()

	a.messages = append(a.messages, message)

	a.mu.Unlock()
}

func (a *API) AddHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	wrappedHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(http.StatusMethodNotAllowed, w, fmt.Sprintf("Method %v not allowed", r.Method))

			return
		}

		w.WriteHeader(http.StatusOK)

		handler(w, r)
	}

	a.mux.HandleFunc(pattern, wrappedHandler)
}

func (a *API) Listen() error {
	err := a.server.ListenAndServe()
	if err != nil && !a.shutdown {
		return err
	}

	return nil
}

func (a *API) Stop() error {
	a.shutdown = true

	return a.server.Shutdown(context.Background()) // FIXME
}
