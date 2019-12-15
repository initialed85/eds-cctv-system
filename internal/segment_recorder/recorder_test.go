package segment_recorder

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

// needs an unauthenticated source at rtsp://192.168.137.31
func TestRecordSegments(t *testing.T) {
	dir, err := ioutil.TempDir("", "file_converter_test")
	if err != nil {
		log.Fatal(err)
	}

	process, err := RecordSegments("rtsp://192.168.137.31", dir, "Driveway", 5)
	if process == nil {
		log.Fatal("process unexpectedly nil")
	}

	if err != nil {
		_ = process.Kill()

		log.Fatal(err)
	}

	time.Sleep(time.Second * 10)

	_ = process.Kill()

	processState, err := process.Wait()
	if err != nil {
		_ = process.Kill()

		log.Fatal(err)
	}

	assert.Greater(t, processState.Pid(), 0)

	assert.Equal(t, -1, processState.ExitCode())

	// TODO: way more with this test
}
