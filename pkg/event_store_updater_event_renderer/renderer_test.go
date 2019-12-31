package event_store_updater_event_renderer

import (
	"github.com/google/uuid"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func getFolderContents(path string) ([]string, error) {
	contents := make([]string, 0)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		_, file := filepath.Split(path)

		contents = append(contents, file)

		return nil
	}

	err := filepath.Walk(path, walkFn)
	if err != nil {
		return []string{}, err
	}

	return contents, nil
}

func TestRenderer(t *testing.T) {
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

	log.Print(dir)

	path := filepath.Join(dir, "some_file.jsonl")

	store := event_store.NewStore(path)
	err = store.Write()
	if err != nil {
		log.Fatal(err)
	}

	renderer, err := New(path, dir)
	if err != nil {
		log.Fatal(err)
	}

	renderer.Start()

	time.Sleep(time.Millisecond * 200)
	folderContents, err := getFolderContents(dir)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(
		t,
		[]string{
			"some_file.jsonl",
		},
		folderContents,
	)

	store.Add(event1)
	err = store.Write()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 200)
	folderContents, err = getFolderContents(dir)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(
		t,
		[]string{
			"events.html",
			"events_0001_01_01.html",
			"some_file.jsonl",
		},
		folderContents,
	)

	time.Sleep(time.Millisecond * 200)
	renderer.Stop()
}
