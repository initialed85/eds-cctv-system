package motion_config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestFind(t *testing.T) {
	configs, err := Find("../../motion-configs")
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t,
		[]Config{
			{CameraName: "Driveway", CameraId: 101, NetCamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/", Width: 0, Height: 0},
			{CameraName: "FrontDoor", CameraId: 102, NetCamURL: "rtsp://192.168.137.32:554/Streaming/Channels/101/", Width: 0, Height: 0},
			{CameraName: "SideGate", CameraId: 103, NetCamURL: "rtsp://192.168.137.33:554/Streaming/Channels/101/", Width: 0, Height: 0},
			{CameraName: "", CameraId: 0, NetCamURL: "", Width: 1920, Height: 1080},
		},
		configs,
	)
}
