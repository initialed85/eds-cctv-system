package file_converter

import (
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/common"
	"log"
)

// ffmpeg -y -i input.mp4 -s 640x360 output.mp4
func ConvertVideo(sourcePath, destinationPath string, width, height int) (string, string, error) {
	log.Printf("converting video %v to %vx%v at %v", sourcePath, width, height, destinationPath)

	arguments := make([]string, 0)

	//if runtime.GOOS != "darwin" {
	//	arguments = append(arguments, "-f")
	//}

	//arguments = append(
	//	arguments,
	//	[]string{
	//		"-l",
	//		"75",
	//		"--",
	//		"ffmpeg",
	//		"-hwaccel",
	//		"cuda",
	//		"-c:v",
	//		"h264_cuvid",
	//		"-y",
	//		"-i",
	//		sourcePath,
	//		"-s",
	//		fmt.Sprintf("%vx%v", width, height),
	//		destinationPath,
	//	}...,
	//)

	arguments = []string{
		"-hwaccel",
		"cuda",
		"-c:v",
		"h264_cuvid",
		"-y",
		"-i",
		sourcePath,
		"-s",
		fmt.Sprintf("%vx%v", width, height),
		destinationPath,
	}

	//return common.RunCommand(
	//	"cpulimit",
	//	arguments...,
	//)

	return common.RunCommand(
		"ffmpeg",
		arguments...,
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
