package motion_log_event_handler

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var lines = []string{
	"[2:ml2:FrontDoor] [NTC] [EVT] event_newfile: File of type 8 saved to: ../../test_files/34__103__2019-12-15_13-38-29__SideGate.mkv",
	"[2:ml2:FrontDoor] [NTC] [ALL] motion_detected: Motion detected - starting event 34",
	"[2:ml2:FrontDoor] [NTC] [EVT] event_newfile: File of type 1 saved to: ../../test_files/34__103__2019-12-15_13-38-31__SideGate.jpg",
	"[2:ml2:FrontDoor] [NTC] [ALL] mlp_actions: End of event 34",
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

	w, err := New(path)
	if err != nil {
		log.Fatal(err)
	}

	w.Start()

	time.Sleep(time.Millisecond * 100)

	lastMessage := "initial"

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "initial", lastMessage)

	err = writeToFile(path, "")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "initial", lastMessage)

	err = appendToFile(path, lines[0]+"\n")
	if err != nil {
		log.Fatal(err)
	}

	err = appendToFile(path, lines[1]+"\n")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)

	err = appendToFile(path, lines[2]+"\n")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)

	err = appendToFile(path, lines[3]+"\n")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	err = w.Stop()
	if err != nil {
		log.Fatal(err)
	}

	_ = os.Remove("../../test_files/34__103__2019-12-15_13-38-29__SideGate-lowres.mkv")
	_ = os.Remove("../../test_files/34__103__2019-12-15_13-38-31__SideGate-lowres.jpg")
}
