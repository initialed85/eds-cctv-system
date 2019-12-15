package common

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"time"
)

func RunCommand(executable string, arguments ...string) (string, string, error) {
	cmd := exec.Command(
		executable,
		arguments...,
	)

	log.Printf("running %v in foreground", cmd.Args)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func RunBackgroundProcess(executable string, arguments ...string) (*os.Process, error) {
	cmd := exec.Command(
		executable,
		arguments...,
	)

	log.Printf("running %v in background", cmd.Args)

	var err error

	go func() {
		err = cmd.Start()
	}()

	time.Sleep(time.Second)

	if err != nil {
		return nil, err
	}

	return cmd.Process, nil
}
