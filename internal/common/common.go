package common

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
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

type BackgroundProcess struct {
	Cmd           *exec.Cmd
	stop, stopped bool
}

func (b *BackgroundProcess) Stop() {
	b.stop = true

	if b.Cmd != nil && b.Cmd.Process != nil {
		_ = b.Cmd.Process.Kill()
	}

	for {
		if b.stopped {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func RunBackgroundProcess(executable string, arguments ...string) (process *BackgroundProcess, startErr error) {
	log.Printf("running %v %v in background", executable, arguments)

	process = &BackgroundProcess{
		stop:    false,
		stopped: false,
	}

	go func() {
		for {
			if process.stop {
				log.Printf("stopping %v %v", executable, arguments)

				_ = process.Cmd.Process.Kill()

				process.stopped = true

				log.Printf("stopped %v %v", executable, arguments)

				break
			}

			log.Printf("creating %v %v", executable, arguments)

			process.Cmd = exec.Command(
				executable,
				arguments...,
			)

			log.Printf("starting %v", process.Cmd.Args)

			startErr = process.Cmd.Start()
			if startErr != nil {
				log.Printf("failed to Start because: %v; trying again...", startErr)

				_ = process.Cmd.Process.Kill()

				time.Sleep(time.Second)

				continue
			}

			log.Printf("waiting for %+v", process.Cmd.Process)

			waitErr := process.Cmd.Wait()
			if waitErr != nil {
				log.Printf("failed to Wait because: %v; trying again...", waitErr)

				_ = process.Cmd.Process.Kill()

				time.Sleep(time.Second)

				continue
			}
		}
	}()

	for {
		if process == nil || process.Cmd == nil || process.Cmd.Process == nil {
			time.Sleep(time.Millisecond * 100)

			continue
		}

		break
	}

	return process, nil
}

func GetLowResPath(path string) string {
	extension := filepath.Ext(path)

	parts := strings.Split(path, extension)

	part := strings.Join(parts[0:len(parts)-1], extension)

	return fmt.Sprintf("%v-lowres%v", part, extension)
}
