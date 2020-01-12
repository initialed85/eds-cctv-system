package event_store

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestMarshalJSONLine(t *testing.T) {
	event := Event{
		CameraName:       "camera1",
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
		"{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n",
		jsonLine,
	)
}

func TestWriteJSONLines(t *testing.T) {
	events := []Event{
		{
			CameraName:       "camera1",
			HighResImagePath: "output/events/image.jpg",
			LowResImagePath:  "output/events/image-lowres.jpg.",
			HighResVideoPath: "output/events/video.mkv",
			LowResVideoPath:  "output/events/video-lowres.mkv",
		},
		{
			CameraName:       "camera1",
			HighResImagePath: "output/events/image.jpg",
			LowResImagePath:  "output/events/image-lowres.jpg.",
			HighResVideoPath: "output/events/video.mkv",
			LowResVideoPath:  "output/events/video-lowres.mkv",
		},
	}

	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.txt")

	for i := 0; i < 2; i++ {
		err = WriteJSONLines(events, path)
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
		"{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n",
		jsonLines,
	)
}

func TestAppendJSONLines(t *testing.T) {
	events := []Event{
		{
			CameraName:       "camera1",
			HighResImagePath: "output/events/image.jpg",
			LowResImagePath:  "output/events/image-lowres.jpg.",
			HighResVideoPath: "output/events/video.mkv",
			LowResVideoPath:  "output/events/video-lowres.mkv",
		},
		{
			CameraName:       "camera1",
			HighResImagePath: "output/events/image.jpg",
			LowResImagePath:  "output/events/image-lowres.jpg.",
			HighResVideoPath: "output/events/video.mkv",
			LowResVideoPath:  "output/events/video-lowres.mkv",
		},
	}

	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.txt")

	for i := 0; i < 2; i++ {
		err = AppendJSONLines(events, path)
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
		"{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n",
		jsonLines,
	)
}

func TestUnmarshalJSONLines(t *testing.T) {
	events, err := UnmarshalJSONLines("{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	assert.Equal(
		t,
		[]Event{
			{
				CameraName:       "camera1",
				HighResImagePath: "output/events/image.jpg",
				LowResImagePath:  "output/events/image-lowres.jpg.",
				HighResVideoPath: "output/events/video.mkv",
				LowResVideoPath:  "output/events/video-lowres.mkv",
			},
			{
				CameraName:       "camera1",
				HighResImagePath: "output/events/image.jpg",
				LowResImagePath:  "output/events/image-lowres.jpg.",
				HighResVideoPath: "output/events/video.mkv",
				LowResVideoPath:  "output/events/video-lowres.mkv",
			},
		},
		events,
	)
}

func TestReadJSONLines(t *testing.T) {
	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	path := filepath.Join(dir, "some_file.txt")

	err = ioutil.WriteFile(
		path,
		[]byte("{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n{\"event_id\":\"00000000-0000-0000-0000-000000000000\",\"timestamp\":\"0001-01-01T00:00:00Z\",\"camera_name\":\"camera1\",\"high_res_image_path\":\"output/events/image.jpg\",\"low_res_image_path\":\"output/events/image-lowres.jpg.\",\"high_res_video_path\":\"output/events/video.mkv\",\"low_res_video_path\":\"output/events/video-lowres.mkv\"}\n"),
		0644,
	)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	events, err := ReadJSONLines(path)

	assert.Equal(
		t,
		[]Event{
			{
				CameraName:       "camera1",
				HighResImagePath: "output/events/image.jpg",
				LowResImagePath:  "output/events/image-lowres.jpg.",
				HighResVideoPath: "output/events/video.mkv",
				LowResVideoPath:  "output/events/video-lowres.mkv",
			},
			{
				CameraName:       "camera1",
				HighResImagePath: "output/events/image.jpg",
				LowResImagePath:  "output/events/image-lowres.jpg.",
				HighResVideoPath: "output/events/video.mkv",
				LowResVideoPath:  "output/events/video-lowres.mkv",
			},
		},
		events,
	)
}
