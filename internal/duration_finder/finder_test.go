package duration_finder

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestFindDuration(t *testing.T) {
	duration, err := FindDuration("../../test_files/34__103__2019-12-15_13-38-29__SideGate.mkv")
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, time.Duration(8260000000), duration)
}
