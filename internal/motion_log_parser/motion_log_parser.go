package motion_log_parser

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Timestamp         time.Time `json:"timestamp"`
	EventId           uuid.UUID `json:"event_id"`
	EventType         string    `json:"event_type"`
	CameraNumber      int64     `json:"camera_number"`
	CameraName        string    `json:"camera_name"`
	CameraEventNumber int64     `json:"camera_event_number"`
}

func GetCameraNumber(line string) (int64, error) {
	cameraNumber, err := strconv.ParseInt(strings.Split(strings.Split(strings.Split(line, "] [", )[0], ":", )[0], "[")[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return cameraNumber, nil
}

func GetCameraName(line string) string {
	return strings.TrimSpace(strings.Split(strings.Split(line, "] [", )[0], ":", )[2])
}

func GetCameraEventNumberFromEndOfLine(line string) (int64, error) {
	parts := strings.Split(line, " ")

	cameraEventNumber, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err != nil {
		return 0, err
	}

	return cameraEventNumber, nil
}

func ParseLine(line string) (Event, error) {
	eventType := "unknown"

	cameraNumber, err := GetCameraNumber(line)
	if err != nil {
		return Event{}, err
	}

	cameraName := GetCameraName(line)

	var cameraEventNumber int64

	if strings.Contains(line, "motion_detected: Motion detected - starting event") {
		eventType = "motion_start"
	} else if strings.Contains(line, "mlp_actions: End of event") {
		eventType = "motion_stop"
	} else if strings.Contains(line, "File of type 8 saved to") {
		eventType = "save_video"
	} else if strings.Contains(line, "File of type 1 saved to") {
		eventType = "save_image"
	}

	if eventType == "motion_start" || eventType == "motion_stop" {
		cameraEventNumber, err = GetCameraEventNumberFromEndOfLine(line)
		if err != nil {
			return Event{}, err
		}
	} else if eventType == "save_video" || eventType == "save_image" {
		cameraEventNumber, err = GetCameraEventNumberFromFilePath(line)
		if err != nil {
			return Event{}, err
		}
	}

	if eventType == "unknown" {
		return Event{}, fmt.Errorf("unsupported event type")
	}

	return Event{
		CameraNumber:      cameraNumber,
		CameraName:        cameraName,
		CameraEventNumber: cameraEventNumber,
	}, nil
}
