package motion_log_parser

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

var lines = []string{
	"[2:ml2:FrontDoor] [NTC] [EVT] event_newfile: File of type 8 saved to: /srv/target_dir/754__102__2019-12-11_13-42-18__FrontDoor.mkv",
	"[2:ml2:FrontDoor] [NTC] [ALL] motion_detected: Motion detected - starting event 754",
	"[2:ml2:FrontDoor] [NTC] [EVT] event_newfile: File of type 1 saved to: /srv/target_dir/754__102__2019-12-11_13-42-19__FrontDoor.jpg",
	"[2:ml2:FrontDoor] [NTC] [ALL] mlp_actions: End of event 754",
}

var nonFileLines = []string{
	lines[1],
	lines[3],
}

var fileLines = []string{
	lines[0],
	lines[2],
}

func TestGetCameraNumber(t *testing.T) {
	for _, line := range lines {
		cameraNumber, err := GetCameraNumber(line)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, int64(2), cameraNumber)
	}
}

func TestGetCameraName(t *testing.T) {
	for _, line := range lines {
		assert.Equal(t, "FrontDoor", GetCameraName(line))
	}
}

func TestGetCameraEventNumberFromEndOfLine(t *testing.T) {
	for _, line := range nonFileLines {
		cameraNumber, err := GetCameraEventNumberFromEndOfLine(line)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, int64(754), cameraNumber)
	}
}

func TestGetCameraEventNumberFromFilePath(t *testing.T) {
	for _, line := range fileLines {
		cameraNumber, err := GetCameraEventNumberFromFilePath(line)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, int64(754), cameraNumber)
	}
}

func TestGetEventID(t *testing.T) {
	assert.Equal(
		t,
		uuid.UUID(uuid.UUID{0xe2, 0x85, 0xed, 0x94, 0x32, 0x8e, 0x52, 0x59, 0xba, 0xb2, 0xb1, 0x55, 0x3b, 0x0, 0xae, 0x5a}),
		GetEventID(2, "FrontDoor", 754),
	)
}

func TestGetFilePath(t *testing.T) {
	assert.Equal(
		t,
		"/srv/target_dir/754__102__2019-12-11_13-42-18__FrontDoor.mkv",
		GetFilePath(lines[0]),
	)

	assert.Equal(
		t,
		"/srv/target_dir/754__102__2019-12-11_13-42-19__FrontDoor.jpg",
		GetFilePath(lines[2]),
	)
}

func TestParseLine(t *testing.T) {
	timestamp := time.Now()

	event, err := ParseLine(timestamp, lines[0])
	if err != nil {
		log.Fatal((err))
	}

	assert.Equal(
		t,
		Event{
			Timestamp:         timestamp,
			EventID:           uuid.UUID{0xe2, 0x85, 0xed, 0x94, 0x32, 0x8e, 0x52, 0x59, 0xba, 0xb2, 0xb1, 0x55, 0x3b, 0x0, 0xae, 0x5a},
			EventState:        "save_video",
			CameraNumber:      2,
			CameraName:        "FrontDoor",
			CameraEventNumber: 754,
			FilePath:          "/srv/target_dir/754__102__2019-12-11_13-42-18__FrontDoor.mkv",
		},
		event,
	)

	event, err = ParseLine(timestamp, lines[1])
	if err != nil {
		log.Fatal((err))
	}

	assert.Equal(
		t,
		Event{
			Timestamp:         timestamp,
			EventID:           uuid.UUID{0xe2, 0x85, 0xed, 0x94, 0x32, 0x8e, 0x52, 0x59, 0xba, 0xb2, 0xb1, 0x55, 0x3b, 0x0, 0xae, 0x5a},
			EventState:        "motion_start",
			CameraNumber:      2,
			CameraName:        "FrontDoor",
			CameraEventNumber: 754,
		},
		event,
	)

	event, err = ParseLine(timestamp, lines[2])
	if err != nil {
		log.Fatal((err))
	}

	assert.Equal(
		t,
		Event{
			Timestamp:         timestamp,
			EventID:           uuid.UUID{0xe2, 0x85, 0xed, 0x94, 0x32, 0x8e, 0x52, 0x59, 0xba, 0xb2, 0xb1, 0x55, 0x3b, 0x0, 0xae, 0x5a},
			EventState:        "save_image",
			CameraNumber:      2,
			CameraName:        "FrontDoor",
			CameraEventNumber: 754,
			FilePath:          "/srv/target_dir/754__102__2019-12-11_13-42-19__FrontDoor.jpg",
		},
		event,
	)

	event, err = ParseLine(timestamp, lines[3])
	if err != nil {
		log.Fatal((err))
	}

	assert.Equal(
		t,
		Event{
			Timestamp:         timestamp,
			EventID:           uuid.UUID{0xe2, 0x85, 0xed, 0x94, 0x32, 0x8e, 0x52, 0x59, 0xba, 0xb2, 0xb1, 0x55, 0x3b, 0x0, 0xae, 0x5a},
			EventState:        "motion_stop",
			CameraNumber:      2,
			CameraName:        "FrontDoor",
			CameraEventNumber: 754,
		},
		event,
	)
}
