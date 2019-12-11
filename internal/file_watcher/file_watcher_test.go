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
	dir, err := ioutil.TempDir("", "watcher_test")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.txt")

	lastAdded := "initial"

	callback := func(added string) {
		lastAdded = added
	}

	w, err := New(path, callback)
	if err != nil {
		log.Fatal(err)
	}

	go w.Watch()

	time.Sleep(time.Millisecond)
	assert.Equal(t, "initial", lastAdded) // nothing yet

	err = writeToFile(path, "")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond)
	assert.Equal(t, "", lastAdded) // file created

	err = appendToFile(path, "The first line\nThe second line\n")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond)
	assert.Equal(t, "The first line\nThe second line\n", lastAdded) // edit w/ trailing newline

	err = appendToFile(path, "The first line\nThe second line")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond)
	assert.Equal(t, "The first line\nThe second line", lastAdded) // edit w/o trailing newline

	err = os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond)
	assert.Equal(t, "The first line\nThe second line", lastAdded) // no change

	w.Stop()
}
