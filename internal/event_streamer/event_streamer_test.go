package event_streamer

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestEventStreamer(t *testing.T) {
	eventStreamer := New(8080)

	go func() {
		err := eventStreamer.Listen()
		if err != nil {
		}
	}()

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/events", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = c.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	lastMessage := "initial"

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			lastMessage = string(message)
		}
	}()

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "initial", lastMessage)

	eventStreamer.AddMessage("Okay, a new message")

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "Okay, a new message", lastMessage)

	eventStreamer.AddMessage("Okay, a even newer message\nOn multiple lines")

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "Okay, a even newer message\nOn multiple lines", lastMessage)
}
