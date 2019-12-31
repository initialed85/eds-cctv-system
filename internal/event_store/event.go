package event_store

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	EventID          uuid.UUID `json:"event_id"`
	Timestamp        time.Time `json:"timestamp"`
	CameraName       string    `json:"camera_name"`
	HighResImagePath string    `json:"high_res_image_path"`
	LowResImagePath  string    `json:"low_res_image_path"`
	HighResVideoPath string    `json:"high_res_video_path"`
	LowResVideoPath  string    `json:"low_res_video_path"`
}

func NewEvent(timestamp time.Time, cameraName, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath string) Event {
	return Event{
		EventID: uuid.NewSHA1(
			uuid.NameSpaceDNS,
			[]byte(fmt.Sprintf("%v, %v, %v, %v, %v", timestamp, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath)),
		),
		Timestamp:        timestamp,
		CameraName:       cameraName,
		HighResImagePath: highResImagePath,
		LowResImagePath:  lowResImagePath,
		HighResVideoPath: highResVideoPath,
		LowResVideoPath:  lowResVideoPath,
	}
}
