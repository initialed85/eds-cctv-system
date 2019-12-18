package folder_watcher

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFileWithFlags(path, data string, flag int) error {
	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(path, data string) error {
	return writeFileWithFlags(path, data, os.O_CREATE|os.O_WRONLY)
}

func appendToFile(path, data string) error {
	return writeFileWithFlags(path, data, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
}

func TestFolderWatcher(t *testing.T) {
	dir, err := ioutil.TempDir("", "folder_watcher_test")
	if err != nil {
		log.Fatal(err)
	}

	path1 := filepath.Join(dir, "Segment_2019-12-15_22-04-28_Driveway.mp4")
	path2 := filepath.Join(dir, "Segment_2019-12-15_22-04-58_Driveway.mp4")
	path3 := filepath.Join(dir, "Segment_2019-12-15_22-05-28_Driveway.mp4")

	lastPath := "initial"

	callback := func(timestamp time.Time, path string) {
		lastPath = path
	}

	f, err := New(dir, callback)
	if err != nil {
		log.Fatal(err)
	}

	go f.Watch()

	time.Sleep(time.Second)
	assert.Equal(t, "initial", lastPath)

	err = writeToFile(path1, "some data")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)
	assert.Equal(t, "initial", lastPath)

	err = writeToFile(path2, "some data")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)
	assert.Equal(t, path1, lastPath)

	err = writeToFile(path3, "some data")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)
	assert.Equal(t, path2, lastPath)
}
