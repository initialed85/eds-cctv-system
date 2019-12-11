package event_streamer

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type EventStreamer struct {
	port     int
	messages []string
	mu       sync.Mutex
	ug       websocket.Upgrader
}

func checkOrigin(r *http.Request) bool {
	return true
}

func New(port int) EventStreamer {
	return EventStreamer{
		port:     port,
		messages: make([]string, 0),
		ug: websocket.Upgrader{
			CheckOrigin: checkOrigin,
		},
	}
}

func (e *EventStreamer) handler(w http.ResponseWriter, r *http.Request) {
	c, err := e.ug.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)

		return
	}

	defer func() {
		err = c.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for {
		if len(e.messages) == 0 {
			time.Sleep(time.Millisecond)
		}

		e.mu.Lock()

		for _, message := range e.messages {
			err = c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println(err)

				break
			}
		}

		e.messages = nil

		e.mu.Unlock()
	}
}

func (e *EventStreamer) Listen() error {
	http.HandleFunc("/events", e.handler)

	return http.ListenAndServe(fmt.Sprintf(":%v", e.port), nil)
}

func (e *EventStreamer) AddMessage(message string) {
	e.mu.Lock()

	e.messages = append(e.messages, message)

	e.mu.Unlock()
}
