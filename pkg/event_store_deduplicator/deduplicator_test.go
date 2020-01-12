package event_store_deduplicator

import (
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
	"time"
)

func TestDeduplicator(t *testing.T) {
	time1 := time.Time{}
	time2 := time1.Add(time.Hour * 24)
	time3 := time2.Add(time.Hour * 24)

	event1 := event_store.NewEvent(time1, "camera1", "image1-hi", "image1-lo", "video1-hi", "video1-lo")
	event2 := event_store.NewEvent(time2, "camera2", "image2-hi", "image2-lo", "video2-hi", "video2-lo")
	event3 := event_store.NewEvent(time3, "camera3", "image3-hi", "image3-lo", "video3-hi", "video3-lo")
	event4 := event_store.NewEvent(event3.Timestamp, "camera3", event3.HighResImagePath, event3.LowResImagePath, event3.HighResVideoPath, event3.LowResVideoPath)

	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	path := filepath.Join(dir, "some_file.txt")

	store := event_store.NewStore(path)
	for i := 0; i < 4; i++ {
		store.Add(event1)
		store.Add(event2)
		store.Add(event3)
		store.Add(event4)

		err := store.Append()
		if err != nil {
			log.Fatal(err)
		}
	}

	assert.Equal(t, 3, store.Len())

	jsonLines, err := event_store.ReadJSONLines(path)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 12, len(jsonLines))

	deduplicator := New(path)

	err = deduplicator.Deduplicate()
	if err != nil {
		log.Fatal(err)
	}

	err = store.Read()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 3, store.Len())

	jsonLines, err = event_store.ReadJSONLines(path)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 3, len(jsonLines))
}