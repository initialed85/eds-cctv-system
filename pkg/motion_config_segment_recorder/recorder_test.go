package motion_config_segment_recorder

import (
	"log"
	"testing"
	"time"
)

func TestMotionConfigSegmentRecorder(t *testing.T) {
	motionConfigSegmentRecorder, err := New("../../motion-configs", "/tmp/", 5)
	if err != nil {
		log.Fatal(err)
	}

	err = motionConfigSegmentRecorder.Start()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 10)

	motionConfigSegmentRecorder.Stop()

	// TODO: way more with this test
}
