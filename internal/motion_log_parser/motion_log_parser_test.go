package motion_log_parser

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
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
