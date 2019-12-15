package common

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

// needs a Linux system with "echo"
func TestRunCommand(t *testing.T) {
	stdout, stderr, err := RunCommand("echo", "hello")
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, "hello\n", stdout)
	assert.Equal(t, "", stderr)
}

// needs a Linux system with "sleep"
func TestRunBackgroundProcess(t *testing.T) {
	before := time.Now()

	process, err := RunBackgroundProcess("sleep", "5")
	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, process)

	processState, err := process.Wait()
	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, processState)

	after := time.Now()

	duration := after.Sub(before)

	assert.Equal(t, true, processState.Exited())

	assert.Greater(t, duration.Seconds(), 5.0)
}
