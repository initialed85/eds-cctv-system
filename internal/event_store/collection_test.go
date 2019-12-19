package event_store

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	time1 := time.Time{}
	time2 := time1.Add(time.Second)
	time3 := time2.Add(time.Second)

	event1 := NewEvent(time1, "image1-hi", "image1-lo", "video1-hi", "video1-lo")
	event2 := NewEvent(time2, "image2-hi", "image2-lo", "video2-hi", "video2-lo")
	event3 := NewEvent(time3, "image3-hi", "image3-lo", "video3-hi", "video3-lo")
	event4 := NewEvent(event3.Timestamp, event3.HighResImagePath, event3.LowResImagePath, event3.HighResVideoPath, event3.LowResVideoPath)

	dir, err := ioutil.TempDir("", "watcher_test")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	path := filepath.Join(dir, "some_file.txt")

	collection := NewStore(path)

	assert.Equal(t, 0, collection.Len())
	assert.Equal(t, []Event{}, collection.GetAll())

	collection.Add(event1)
	collection.Add(event2)
	collection.Add(event3)
	collection.Add(event3)
	collection.Add(event4)

	assert.Equal(t, 3, collection.Len())
	assert.Equal(
		t,
		[]Event{
			{EventID: uuid.UUID{0x6, 0xb1, 0xdd, 0x33, 0x71, 0xf3, 0x58, 0x16, 0xb0, 0x92, 0x8b, 0x87, 0x6a, 0xb8, 0xfc, 0xc6}, Timestamp: time1, HighResImagePath: "image1-hi", LowResImagePath: "image1-lo", HighResVideoPath: "video1-hi", LowResVideoPath: "video1-lo"},
			{EventID: uuid.UUID{0x88, 0x4f, 0x1b, 0xc1, 0x8c, 0x12, 0x5d, 0x77, 0xac, 0xa6, 0x72, 0x8, 0xa, 0x65, 0x96, 0xe2}, Timestamp: time2, HighResImagePath: "image2-hi", LowResImagePath: "image2-lo", HighResVideoPath: "video2-hi", LowResVideoPath: "video2-lo"},
			{EventID: uuid.UUID{0x13, 0x2, 0xdb, 0x5e, 0xd4, 0xde, 0x5f, 0x73, 0x81, 0x79, 0xec, 0xb8, 0xd2, 0xa3, 0x78, 0x4a}, Timestamp: time3, HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"},
		},
		collection.GetAll(),
	)

	err = collection.Remove(event2.EventID)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	assert.Equal(t, 2, collection.Len())
	assert.Equal(
		t,
		[]Event{
			{EventID: uuid.UUID{0x6, 0xb1, 0xdd, 0x33, 0x71, 0xf3, 0x58, 0x16, 0xb0, 0x92, 0x8b, 0x87, 0x6a, 0xb8, 0xfc, 0xc6}, Timestamp: time1, HighResImagePath: "image1-hi", LowResImagePath: "image1-lo", HighResVideoPath: "video1-hi", LowResVideoPath: "video1-lo"},
			{EventID: uuid.UUID{0x13, 0x2, 0xdb, 0x5e, 0xd4, 0xde, 0x5f, 0x73, 0x81, 0x79, 0xec, 0xb8, 0xd2, 0xa3, 0x78, 0x4a}, Timestamp: time3, HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"},
		},
		collection.GetAll(),
	)

	event, err := collection.Pop(event1.EventID)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	assert.Equal(t, event1, event)
	assert.Equal(t, 1, collection.Len())
	assert.Equal(
		t,
		[]Event{
			{EventID: uuid.UUID{0x13, 0x2, 0xdb, 0x5e, 0xd4, 0xde, 0x5f, 0x73, 0x81, 0x79, 0xec, 0xb8, 0xd2, 0xa3, 0x78, 0x4a}, Timestamp: time3, HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"},
		},
		collection.GetAll(),
	)

	event, err = collection.Get(event3.EventID)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	assert.Equal(t, event3, event)
	assert.Equal(t, 1, collection.Len())
	assert.Equal(
		t,
		[]Event{
			{EventID: uuid.UUID{0x13, 0x2, 0xdb, 0x5e, 0xd4, 0xde, 0x5f, 0x73, 0x81, 0x79, 0xec, 0xb8, 0xd2, 0xa3, 0x78, 0x4a}, Timestamp: time3, HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"},
		},
		collection.GetAll(),
	)
}
