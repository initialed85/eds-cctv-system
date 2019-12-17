package file_converter

import (
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/common"
	"log"
)

// ffmpeg -i input.mp4 -s 640x480 output.mp4
func ConvertVideo(sourcePath, destinationPath string, width, height int) (string, string, error) {
	log.Printf("converting video %v to %vx%v at %v", sourcePath, width, height, destinationPath)

	return common.RunCommand(
		"ffmpeg",
		"-i",
		sourcePath,
		"-s",
		fmt.Sprintf("%vx%v", width, height),
		destinationPath,
	)
}

// convert -resize 1024X768 source.png dest.jpg
func ConvertImage(sourcePath, destinationPath string, width, height int) (string, string, error) {
	log.Printf("converting image %v to %vx%v at %v", sourcePath, width, height, destinationPath)

	return common.RunCommand(
		"convert",
		"-resize",
		fmt.Sprintf("%vX%v", width, height),
		sourcePath,
		destinationPath,
	)
}
