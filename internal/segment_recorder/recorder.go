package segment_recorder

import (
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/common"
	"log"
	"path/filepath"
)

// ffmpeg -rtsp_transport tcp -i rtsp://192.168.137.31:554/Streaming/Channels/101 -c copy -map 0 -f segment -segment_time 60 -segment_format mp4 -segment_atclocktime 1 -strftime 1 -g 10 /srv/target_dir/segments/Segment_%Y-%m-%d_%H-%M-%S_Driveway.mp4
func RecordSegments(netCamURL, destinationPath, cameraName string, duration int) (*common.BackgroundProcess, error) {
	log.Printf("recording %v second segments from %v to %v for %v", duration, netCamURL, destinationPath, cameraName)

	return common.RunBackgroundProcess(
		"ffmpeg",
		"-hwaccel",
		"vaapi",
		"-rtsp_transport",
		"tcp",
		"-i",
		netCamURL,
		"-c",
		"copy",
		"-map",
		"0",
		"-f",
		"segment",
		"-segment_time",
		fmt.Sprintf("%v", duration),
		"-segment_format",
		"mp4",
		"-segment_atclocktime",
		"1",
		"-strftime",
		"1",
		"-x264-params",
		"keyint=100:scenecut=0",
		"-g",
		"100",
		"-muxdelay",
		"0",
		"-muxpreload",
		"0",
		"-reset_timestamps",
		" 1",
		filepath.Join(destinationPath, "Segment_%Y-%m-%d_%H-%M-%S_"+cameraName+".mp4"),
	)
}
