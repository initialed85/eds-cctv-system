package motion_config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestParseFile(t *testing.T) {
	config, err := ParseFile("../../motion-configs/motion.conf")
	if err != nil {
		log.Fatal("during test:", err)
	}

	assert.Equal(
		t,
		Config{CameraName: "", CameraId: 0, NetCamURL: "", Width: 1920, Height: 1080},
		config,
	)

	config, err = ParseFile("../../motion-configs/conf.d/camera_Driveway.conf")
	if err != nil {
		log.Fatal("during test:", err)
	}

	assert.Equal(
		t,
		Config{CameraName: "Driveway", CameraId: 101, NetCamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/", Width: 0, Height: 0},
		config,
	)
}

func TestMergeConfigs(t *testing.T) {
	leftConfig, err := ParseFile("../../motion-configs/motion.conf")
	if err != nil {
		log.Fatal("during test:", err)
	}

	rightConfig, err := ParseFile("../../motion-configs/conf.d/camera_Driveway.conf")
	if err != nil {
		log.Fatal("during test:", err)
	}

	config := MergeConfigs(leftConfig, rightConfig)

	assert.Equal(
		t,
		Config{CameraName: "Driveway", CameraId: 101, NetCamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/", Width: 1920, Height: 1080},
		config,
	)
}
