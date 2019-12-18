package event_store

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
	"time"
)

func TestEventCollection(t *testing.T) {
	event1 := NewEvent(time.Time{}, "image1-hi", "image1-lo", "video1-hi", "video1-lo")
	event2 := NewEvent(time.Time{}, "image2-hi", "image2-lo", "video2-hi", "video2-lo")
	event3 := NewEvent(time.Time{}, "image3-hi", "image3-lo", "video3-hi", "video3-lo")
	event4 := NewEvent(event3.Timestamp, event3.HighResImagePath, event3.LowResImagePath, event3.HighResVideoPath, event3.LowResVideoPath)

	dir, err := ioutil.TempDir("", "watcher_test")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	path := filepath.Join(dir, "some_file.txt")

	collection := NewEventCollection(path)

	assert.Equal(t, 0, collection.Len())

	collection.Add(event1)
	collection.Add(event2)
	collection.Add(event3)
	collection.Add(event3)
	collection.Add(event4)

	log.Printf("%+v", event1)
	log.Printf("%+v", event2)
	log.Printf("%+v", event3)
	log.Printf("%+v", event4)

	assert.Equal(t, 3, collection.Len())
}
