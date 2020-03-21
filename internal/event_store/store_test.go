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
	time2 := time1.Add(time.Hour * 24)
	time3 := time2.Add(time.Hour * 24)
	time4 := time3.Add(time.Hour * 24)

	event1 := NewEvent(time1, "camera1", "image1-hi", "image1-lo", "video1-hi", "video1-lo")
	event2 := NewEvent(time2, "camera2", "image2-hi", "image2-lo", "video2-hi", "video2-lo")
	event3 := NewEvent(time3, "camera3", "image3-hi", "image3-lo", "video3-hi", "video3-lo")
	event4 := NewEvent(event3.Timestamp, "camera3", event3.HighResImagePath, event3.LowResImagePath, event3.HighResVideoPath, event3.LowResVideoPath)
	event5 := NewEvent(time4, "camera3", event3.HighResImagePath, event3.LowResImagePath, event3.HighResVideoPath, event3.LowResVideoPath)

	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	path := filepath.Join(dir, "some_file.txt")

	store := NewStore(path)

	assert.Equal(t, 0, store.Len())
	assert.Equal(t, []Event{}, store.GetAllDescending())

	store.Overwrite(event1)
	store.Overwrite(event1)
	store.Overwrite(event2)
	store.Overwrite(event2)
	store.Overwrite(event3)
	store.Overwrite(event3)
	store.Overwrite(event4)
	store.Overwrite(event4)

	assert.Equal(t, 3, store.Len())
	assert.Equal(
		t,
		[]Event{
			{EventID: uuid.UUID{0xa2, 0xea, 0xf2, 0xdf, 0xff, 0x28, 0x58, 0x1b, 0xac, 0xf3, 0x5a, 0x63, 0x8d, 0xec, 0xa7, 0xd}, Timestamp: time3, CameraName: "camera3", HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"},
			{EventID: uuid.UUID{0xb8, 0x8, 0x6, 0x48, 0x50, 0xfc, 0x54, 0xb6, 0x83, 0xd6, 0x73, 0x28, 0xb5, 0x13, 0x3c, 0x8f}, Timestamp: time2, CameraName: "camera2", HighResImagePath: "image2-hi", LowResImagePath: "image2-lo", HighResVideoPath: "video2-hi", LowResVideoPath: "video2-lo"},
			{EventID: uuid.UUID{0x6, 0xb1, 0xdd, 0x33, 0x71, 0xf3, 0x58, 0x16, 0xb0, 0x92, 0x8b, 0x87, 0x6a, 0xb8, 0xfc, 0xc6}, Timestamp: time1, CameraName: "camera1", HighResImagePath: "image1-hi", LowResImagePath: "image1-lo", HighResVideoPath: "video1-hi", LowResVideoPath: "video1-lo"},
		},
		store.GetAllDescending(),
	)

	assert.Equal(
		t,
		map[time.Time][]Event{
			time3: {{EventID: uuid.UUID{0xa2, 0xea, 0xf2, 0xdf, 0xff, 0x28, 0x58, 0x1b, 0xac, 0xf3, 0x5a, 0x63, 0x8d, 0xec, 0xa7, 0xd}, Timestamp: time3, CameraName: "camera3", HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"}},
			time2: {{EventID: uuid.UUID{0xb8, 0x8, 0x6, 0x48, 0x50, 0xfc, 0x54, 0xb6, 0x83, 0xd6, 0x73, 0x28, 0xb5, 0x13, 0x3c, 0x8f}, Timestamp: time2, CameraName: "camera2", HighResImagePath: "image2-hi", LowResImagePath: "image2-lo", HighResVideoPath: "video2-hi", LowResVideoPath: "video2-lo"}},
			time1: {{EventID: uuid.UUID{0x6, 0xb1, 0xdd, 0x33, 0x71, 0xf3, 0x58, 0x16, 0xb0, 0x92, 0x8b, 0x87, 0x6a, 0xb8, 0xfc, 0xc6}, Timestamp: time1, CameraName: "camera1", HighResImagePath: "image1-hi", LowResImagePath: "image1-lo", HighResVideoPath: "video1-hi", LowResVideoPath: "video1-lo"}},
		},
		store.GetAllDescendingByDateDescending(),
	)

	event, err := store.GetByUUID(event3.EventID)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	assert.Equal(t, event3, event)
	assert.Equal(t, 3, store.Len())
	assert.Equal(
		t,
		[]Event{
			{EventID: uuid.UUID{0xa2, 0xea, 0xf2, 0xdf, 0xff, 0x28, 0x58, 0x1b, 0xac, 0xf3, 0x5a, 0x63, 0x8d, 0xec, 0xa7, 0xd}, Timestamp: time3, CameraName: "camera3", HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"},
			{EventID: uuid.UUID{0xb8, 0x8, 0x6, 0x48, 0x50, 0xfc, 0x54, 0xb6, 0x83, 0xd6, 0x73, 0x28, 0xb5, 0x13, 0x3c, 0x8f}, Timestamp: time2, CameraName: "camera2", HighResImagePath: "image2-hi", LowResImagePath: "image2-lo", HighResVideoPath: "video2-hi", LowResVideoPath: "video2-lo"},
			{EventID: uuid.UUID{0x6, 0xb1, 0xdd, 0x33, 0x71, 0xf3, 0x58, 0x16, 0xb0, 0x92, 0x8b, 0x87, 0x6a, 0xb8, 0xfc, 0xc6}, Timestamp: time1, CameraName: "camera1", HighResImagePath: "image1-hi", LowResImagePath: "image1-lo", HighResVideoPath: "video1-hi", LowResVideoPath: "video1-lo"},
		},
		store.GetAllDescending(),
	)

	store.Add(event5)
	event, err = store.GetByUUID(event5.EventID)
	assert.NotNil(t, err)

	assert.Equal(
		t,
		[]Event{
			{EventID: uuid.UUID{0x6, 0xb1, 0xdd, 0x33, 0x71, 0xf3, 0x58, 0x16, 0xb0, 0x92, 0x8b, 0x87, 0x6a, 0xb8, 0xfc, 0xc6}, Timestamp: time1, CameraName: "camera1", HighResImagePath: "image1-hi", LowResImagePath: "image1-lo", HighResVideoPath: "video1-hi", LowResVideoPath: "video1-lo"},
			{EventID: uuid.UUID{0xb8, 0x8, 0x6, 0x48, 0x50, 0xfc, 0x54, 0xb6, 0x83, 0xd6, 0x73, 0x28, 0xb5, 0x13, 0x3c, 0x8f}, Timestamp: time2, CameraName: "camera2", HighResImagePath: "image2-hi", LowResImagePath: "image2-lo", HighResVideoPath: "video2-hi", LowResVideoPath: "video2-lo"},
			{EventID: uuid.UUID{0xa2, 0xea, 0xf2, 0xdf, 0xff, 0x28, 0x58, 0x1b, 0xac, 0xf3, 0x5a, 0x63, 0x8d, 0xec, 0xa7, 0xd}, Timestamp: time3, CameraName: "camera3", HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"},
		},
		store.GetAll(),
	)
}
