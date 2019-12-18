package file_watcher

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

func TestFileWatcher(t *testing.T) {
	dir, err := ioutil.TempDir("", "file_watcher_test")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	path := filepath.Join(dir, "some_file.txt")

	lastAdded := []string{"initial"}

	callback := func(timestamp time.Time, added []string) {
		lastAdded = added
	}

	w, err := New(path, callback)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	go w.Watch()

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, []string{"initial"}, lastAdded) // nothing yet

	err = writeToFile(path, "")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	time.Sleep(time.Second)
	assert.Equal(t, []string{""}, lastAdded) // file created

	err = appendToFile(path, "The first line\nThe second line\n")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, []string{"The first line", "The second line"}, lastAdded) // edit w/ trailing newline

	err = appendToFile(path, "The first line\nThe second line")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, []string{"The first line", "The second line"}, lastAdded) // edit w/o trailing newline

	err = os.Remove(path)
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, []string{"The first line", "The second line"}, lastAdded) // no change

	w.Stop()
}
