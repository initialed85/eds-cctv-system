package event_store_updater

import (
	"github.com/google/uuid"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
	"time"
)

func TestUpdater(t *testing.T) {
	time1 := time.Time{}
	time2 := time1.Add(time.Hour * 24)
	time3 := time2.Add(time.Hour * 24)

	event1 := event_store.NewEvent(time1, "camera1", "image1-hi", "image1-lo", "video1-hi", "video1-lo")
	event2 := event_store.NewEvent(time2, "camera2", "image2-hi", "image2-lo", "video2-hi", "video2-lo")
	event3 := event_store.NewEvent(time3, "camera3", "image3-hi", "image3-lo", "video3-hi", "video3-lo")

	event1.EventID = uuid.UUID{0}
	event2.EventID = uuid.UUID{1}
	event3.EventID = uuid.UUID{2}

	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	path := filepath.Join(dir, "some_file.jsonl")

	store := event_store.NewStore(path)

	updated := false
	eventsByDate := make(map[time.Time][]event_store.Event)
	callback := func(store event_store.Store) {
		updated = true
		eventsByDate = store.GetAllByDate()
	}

	updater, err := New(store, callback)
	if err != nil {
		log.Fatal(err)
	}

	go updater.Watch()

	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, false, updated)

	updated = false
	store.Add(event1)
	err = store.Append()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, true, updated)
	assert.Equal(
		t,
		map[time.Time][]event_store.Event{
			time1: {{EventID: uuid.UUID{0}, Timestamp: time1, CameraName: "camera1", HighResImagePath: "image1-hi", LowResImagePath: "image1-lo", HighResVideoPath: "video1-hi", LowResVideoPath: "video1-lo"}},
		},
		eventsByDate,
	)

	updated = false
	store.Add(event2)
	store.Add(event3)
	err = store.Append()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, true, updated)
	assert.Equal(
		t,
		map[time.Time][]event_store.Event{
			time1: {{EventID: uuid.UUID{0}, Timestamp: time1, CameraName: "camera1", HighResImagePath: "image1-hi", LowResImagePath: "image1-lo", HighResVideoPath: "video1-hi", LowResVideoPath: "video1-lo"}},
			time2: {{EventID: uuid.UUID{1}, Timestamp: time2, CameraName: "camera2", HighResImagePath: "image2-hi", LowResImagePath: "image2-lo", HighResVideoPath: "video2-hi", LowResVideoPath: "video2-lo"}},
			time3: {{EventID: uuid.UUID{2}, Timestamp: time3, CameraName: "camera3", HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"}},
		},
		eventsByDate,
	)

	time.Sleep(time.Millisecond * 200)
	updater.Stop()
}
