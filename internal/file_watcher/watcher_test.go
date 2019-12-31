package file_watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	eventReceived := false
	lastEvent := fsnotify.Event{}

	callback := func(event fsnotify.Event, time time.Time) bool {
		log.Printf("event=%+v, time=%+v", event, time)

		eventReceived = true
		lastEvent = event

		return true
	}

	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatal(err)
	}

	extraPath := filepath.Join(dir, "some_folder")

	w, err := New(extraPath, callback)
	if err != nil {
		log.Fatal(err)
	}

	go w.Watch()

	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, false, eventReceived)

	err = os.Mkdir(extraPath, 0700)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, false, eventReceived)

	path1 := filepath.Join(extraPath, "some_file")

	file, err := os.Create(path1)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, true, eventReceived)
	assert.Equal(t, fsnotify.Create, lastEvent.Op)
	assert.Equal(t, path1, lastEvent.Name)
	eventReceived = false

	_, err = file.WriteString("some data\n")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, true, eventReceived)
	assert.Equal(t, fsnotify.Write, lastEvent.Op)
	assert.Equal(t, path1, lastEvent.Name)
	eventReceived = false

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, false, eventReceived)

	err = os.Remove(path1)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, true, eventReceived)
	assert.Equal(t, fsnotify.Remove, lastEvent.Op)
	assert.Equal(t, path1, lastEvent.Name)
	w.Stop()
}
