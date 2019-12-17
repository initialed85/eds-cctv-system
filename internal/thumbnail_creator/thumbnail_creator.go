package thumbnail_creator

import (
	"github.com/initialed85/eds-cctv-system/internal/common"
	"log"
)

func CreateThumbnail(videoPath, imagePath string) error {
	_, stderr, err := common.RunCommand(
		"ffmpeg",
		"-i",
		videoPath,
		"-ss",
		"00:00:01.000",
		"-vframes",
		"1",
		imagePath,
	)

	if err != nil {
		log.Fatal(stderr)
	}

	return err
}
