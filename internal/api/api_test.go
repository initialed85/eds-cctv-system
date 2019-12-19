package api

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestAPI_Stream(t *testing.T) {
	api := New(8080)

	go func() {
		err := api.Listen()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/stream", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = c.Close()
		if err != nil {
			log.Fatal("Close during test:", err)
		}
	}()

	lastMessage := "initial"

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Print(err)

				break
			}

			lastMessage = string(message)
		}
	}()

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "initial", lastMessage)

	api.AddMessage("Okay, a new message")

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "Okay, a new message", lastMessage)

	api.AddMessage("Okay, a even newer message\nOn multiple lines")

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "Okay, a even newer message\nOn multiple lines", lastMessage)

	err = api.Stop()
	if err != nil {
		log.Fatal(err)
	}
}

func TestAPI_Request(t *testing.T) {
	api := New(8080)

	go func() {
		err := api.Listen()
		if err != nil {
			log.Fatalf("during test listen: %v", err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello, world.\n")
		if err != nil {
			log.Fatalf("during test handler: %v", err)
		}
	}

	api.AddHandler("/events", handler)

	time.Sleep(time.Millisecond * 100)

	resp, err := http.Get("http://localhost:8080/events")
	if err != nil {
		log.Fatalf("during test get: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("during test readall: %v", err)
	}

	assert.Equal(t, "Hello, world.\n", string(data))
}
