package segment_folder_event_handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestSegmentFolderEventHandler(t *testing.T) {
	defer func() {
		_ = os.Remove("../../test_files/Segment_a_Driveway.mp4")
		_ = os.Remove("../../test_files/Segment_a_Driveway.jpg")
		_ = os.Remove("../../test_files/Segment_a_Driveway-lowres.mp4")
		_ = os.Remove("../../test_files/Segment_a_Driveway-lowres.jpg")

		_ = os.Remove("../../test_files/Segment_b_Driveway.mp4")
		_ = os.Remove("../../test_files/Segment_b_Driveway.jpg")
		_ = os.Remove("../../test_files/Segment_b_Driveway-lowres.mp4")
		_ = os.Remove("../../test_files/Segment_b_Driveway-lowres.jpg")

		_ = os.Remove("../../test_files/Segment_c_Driveway.mp4")
		_ = os.Remove("../../test_files/Segment_c_Driveway.jpg")
		_ = os.Remove("../../test_files/Segment_c_Driveway-lowres.mp4")
		_ = os.Remove("../../test_files/Segment_c_Driveway-lowres.jpg")
	}()

	var lastCameraName string
	var lastHighResImagePath string
	var lastLowResImagePath string
	var lastHighResVideoPath string
	var lastLowResVideoPath string

	callback := func(timestamp time.Time, cameraName, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath string) error {
		fmt.Println(timestamp, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath)

		lastCameraName = cameraName
		lastHighResImagePath = highResImagePath
		lastLowResImagePath = lowResImagePath
		lastHighResVideoPath = highResVideoPath
		lastLowResVideoPath = lowResVideoPath

		return nil
	}

	s, err := New("../../test_files", callback)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	s.Start()

	time.Sleep(time.Second)
	assert.Equal(t, "", lastCameraName)
	assert.Equal(t, "", lastHighResImagePath)
	assert.Equal(t, "", lastLowResImagePath)
	assert.Equal(t, "", lastHighResVideoPath)
	assert.Equal(t, "", lastLowResVideoPath)

	from, err := ioutil.ReadFile("../../test_files/34__103__2019-12-15_13-38-29__SideGate.mkv")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	err = ioutil.WriteFile("../../test_files/Segment_a_Driveway.mp4", from, 0644)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	time.Sleep(time.Second * 5)
	assert.Equal(t, "", lastCameraName)
	assert.Equal(t, "", lastHighResImagePath)
	assert.Equal(t, "", lastLowResImagePath)
	assert.Equal(t, "", lastHighResVideoPath)
	assert.Equal(t, "", lastLowResVideoPath)

	err = ioutil.WriteFile("../../test_files/Segment_b_Driveway.mp4", from, 0644)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	time.Sleep(time.Second * 5)
	assert.Equal(t, "Driveway", lastCameraName)
	assert.Equal(t, "../../test_files/Segment_a_Driveway.jpg", lastHighResImagePath)
	assert.Equal(t, "../../test_files/Segment_a_Driveway-lowres.jpg", lastLowResImagePath)
	assert.Equal(t, "../../test_files/Segment_a_Driveway.mp4", lastHighResVideoPath)
	assert.Equal(t, "../../test_files/Segment_a_Driveway-lowres.mp4", lastLowResVideoPath)

	err = ioutil.WriteFile("../../test_files/Segment_c_Driveway.mp4", from, 0644)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	time.Sleep(time.Second * 5)
	assert.Equal(t, "Driveway", lastCameraName)
	assert.Equal(t, "../../test_files/Segment_b_Driveway.jpg", lastHighResImagePath)
	assert.Equal(t, "../../test_files/Segment_b_Driveway-lowres.jpg", lastLowResImagePath)
	assert.Equal(t, "../../test_files/Segment_b_Driveway.mp4", lastHighResVideoPath)
	assert.Equal(t, "../../test_files/Segment_b_Driveway-lowres.mp4", lastLowResVideoPath)

	err = s.Stop()
	if err != nil {
		log.Printf("during test: %v", err)
	}

	time.Sleep(time.Second)
}
