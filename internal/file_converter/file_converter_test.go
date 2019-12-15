package file_converter

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// needs ffmpeg installed
func TestConvertVideo(t *testing.T) {
	dir, err := ioutil.TempDir("", "file_converter_test")
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
		"../../test_files/34__103__2019-12-15_13-38-29__SideGate.mkv",
		path,
		640,
		480,
	)
	if err != nil {
		log.Fatalf("during test: %v; stderr= %v", err, stderr)
	}

	assert.NotEqual(t, "", stderr)
	assert.Equal(t, "", stdout)

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal("during test:", err)
	}
}

// needs imagemagick installed
func TestConvertImage(t *testing.T) {
	dir, err := ioutil.TempDir("", "file_converter_test")
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
		"../../test_files/34__103__2019-12-15_13-38-31__SideGate.jpg",
		path,
		640,
		480,
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
