package thumbnail_creator

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestCreateThumbnail(t *testing.T) {
	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.jpg")

	err = CreateThumbnail("../../test_files/34__103__2019-12-15_13-38-29__FrontDoor.mkv", path)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: way more with this test
}
