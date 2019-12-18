package event_persistence

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestGetJSONLine(t *testing.T) {
	event := Event{
		HighResImagePath: "output/events/image.jpg",
		LowResImagePath:  "output/events/image-lowres.jpg.",
		HighResVideoPath: "output/events/video.mkv",
		LowResVideoPath:  "output/events/video-lowres.mkv",
	}

	jsonLine, err := MarshalJSONLine(event)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	assert.Equal(
		t,
		"{\"timestamp\":\"0001-01-01T00:00:00Z\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n",
		jsonLine,
	)
}

func TestWriteJSONLine(t *testing.T) {
	event := Event{
		HighResImagePath: "output/events/image.jpg",
		LowResImagePath:  "output/events/image-lowres.jpg.",
		HighResVideoPath: "output/events/video.mkv",
		LowResVideoPath:  "output/events/video-lowres.mkv",
	}

	dir, err := ioutil.TempDir("", "watcher_test")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.txt")

	for i := 0; i < 2; i++ {
		err = WriteJSONLine(event, path)
		if err != nil {
			log.Fatalf("during test: %v", err)
		}
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	jsonLines := string(b)

	assert.Equal(
		t,
		"{\"timestamp\":\"0001-01-01T00:00:00Z\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n{\"timestamp\":\"0001-01-01T00:00:00Z\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n",
		jsonLines,
	)
}
