package event_streamer

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

func checkOrigin(r *http.Request) bool {
	return true
}

type EventStreamer struct {
	mux      *http.ServeMux
	server   http.Server
	ug       websocket.Upgrader
	mu       sync.Mutex
	messages []string
	shutdown bool
}

func New(port int64) *EventStreamer {
	mux := http.NewServeMux()

	e := &EventStreamer{
		mux: mux,
		server: http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: mux,
		},
		ug:       websocket.Upgrader{CheckOrigin: checkOrigin},
		messages: make([]string, 0),
		shutdown: false,
	}

	e.mux.HandleFunc("/events", e.handleEvents)

	return e
}

func (e *EventStreamer) AddMessage(message string) {
	e.mu.Lock()

	e.messages = append(e.messages, message)

	e.mu.Unlock()
}

func (e *EventStreamer) handleEvents(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleEvents called")

	c, err := e.ug.Upgrade(w, r, nil)
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
		if len(e.messages) == 0 {
			time.Sleep(time.Millisecond)

			continue
		}

		log.Printf("sending %v messages", len(e.messages))

		err = nil

		e.mu.Lock()

		for _, message := range e.messages {
			err = c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("bailing on message loop because: %v", err)

				break
			}
		}

		e.messages = nil

		e.mu.Unlock()

		if err != nil {
			log.Printf("bailing on handleEvents because: %v", err)

			break
		}
	}
}

func (e *EventStreamer) Listen() error {
	err := e.server.ListenAndServe()
	if err != nil && !e.shutdown {
		return err
	}

	return nil
}

func (e *EventStreamer) Stop() error {
	e.shutdown = true

	return e.server.Shutdown(context.Background()) // FIXME
}
