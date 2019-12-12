package motion_log_event_streamer

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var lines = []string{
	"[2:ml2:FrontDoor] [NTC] [EVT] event_newfile: File of type 8 saved to: /srv/target_dir/754__102__2019-12-11_13-42-18__FrontDoor.mkv",
	"[2:ml2:FrontDoor] [NTC] [ALL] motion_detected: Motion detected - starting event 754",
	"[2:ml2:FrontDoor] [NTC] [EVT] event_newfile: File of type 1 saved to: /srv/target_dir/754__102__2019-12-11_13-42-19__FrontDoor.jpg",
	"[2:ml2:FrontDoor] [NTC] [ALL] mlp_actions: End of event 754",
}

func writeFileWithFlags(path, data string, flag int) error {
	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(path, data string) error {
	return writeFileWithFlags(path, data, os.O_CREATE|os.O_WRONLY)
}

func appendToFile(path, data string) error {
	return writeFileWithFlags(path, data, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
}

func TestMotionLogEventStreamer(t *testing.T) {
	dir, err := ioutil.TempDir("", "watcher_test")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.txt")

	w, err := New(path, 8080)
	if err != nil {
		log.Fatal(err)
	}

	w.Start()

	time.Sleep(time.Millisecond * 100)

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

	time.Sleep(time.Millisecond * 100)

	lastMessage := "initial"

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Print("ReadMessage during test:", err)

				break
			}

			lastMessage = string(message)
		}
	}()

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "initial", lastMessage)

	err = writeToFile(path, "")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "initial", lastMessage)

	err = appendToFile(path, lines[0])
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)
	assert.Equal(
		t,
		true,
		strings.HasSuffix(lastMessage, "\"event_id\":\"e285ed94-328e-5259-bab2-b1553b00ae5a\",\"event_state\":\"save_video\",\"camera_number\":2,\"camera_name\":\"FrontDoor\",\"camera_event_number\":754,\"file_path\":\"/srv/target_dir/754__102__2019-12-11_13-42-18__FrontDoor.mkv\"}"),
	)

	err = appendToFile(path, lines[1])
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)
	assert.Equal(
		t,
		true,
		strings.HasSuffix(lastMessage, "\"event_id\":\"e285ed94-328e-5259-bab2-b1553b00ae5a\",\"event_state\":\"motion_start\",\"camera_number\":2,\"camera_name\":\"FrontDoor\",\"camera_event_number\":754,\"file_path\":\"\"}"),
	)

	err = w.Stop()
	if err != nil {
		log.Fatal(err)
	}
}
