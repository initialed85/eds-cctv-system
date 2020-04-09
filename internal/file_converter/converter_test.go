package file_converter

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// needs ffmpeg and cpulimit installed
func TestConvertVideo(t *testing.T) {
	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.mkv")

	defer func() {
		err = os.Remove(path)
		if err != nil {
			log.Print("defer during test:", err)
		}
	}()

	stdout, stderr, err := ConvertVideo(
		"../../test_files/34__103__2019-12-15_13-38-29__FrontDoor.mkv",
		path,
		640,
		360,
	)
	if err != nil {
		log.Fatalf("during test: %v; stderr= %v", err, stderr)
	}

	assert.NotEqual(t, "", stderr)
	if runtime.GOOS == "darwin" {
		assert.Equal(t, "", stdout)
	} else {
		assert.NotEqual(t, "", stdout)
	}

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal("during test:", err)
	}
}

// needs imagemagick installed
func TestConvertImage(t *testing.T) {
	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.jpg")

	defer func() {
		err = os.Remove(path)
		if err != nil {
			log.Print("defer during test:", err)
		}
	}()

	stdout, stderr, err := ConvertImage(
		"../../test_files/34__103__2019-12-15_13-38-31__FrontDoor.jpg",
		path,
		640,
		360,
	)
	if err != nil {
		log.Fatalf("during test: %v; stderr= %v", err, stderr)
	}

	assert.Equal(t, "", stderr)
	assert.Equal(t, "", stdout)

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal("during test:", err)
	}
}
