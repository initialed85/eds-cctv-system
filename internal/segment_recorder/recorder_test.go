package segment_recorder

import (
	"io/ioutil"
	"log"
	"testing"
	"time"
)

// needs an unauthenticated source at rtsp://192.168.137.31
func TestRecordSegments(t *testing.T) {
	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatal(err)
	}

	process, err := RecordSegments("rtsp://192.168.137.31", dir, "Driveway", 5)
	if process == nil {
		log.Fatal("process unexpectedly nil")
	}

	if err != nil {
		process.Stop()

		log.Fatal(err)
	}

	time.Sleep(time.Second * 10)

	process.Stop()

	// TODO: way more with this test
}
