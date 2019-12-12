package motion_log_parser

import (
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Timestamp         time.Time `json:"timestamp"`
	EventID           uuid.UUID `json:"event_id"`
	EventState        string    `json:"event_state"`
	CameraNumber      int64     `json:"camera_number"`
	CameraName        string    `json:"camera_name"`
	CameraEventNumber int64     `json:"camera_event_number"`
	FilePath          string    `json:"file_path"`
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

// NOTE: Relies on filename being prefixed with "(event number)__"
func GetCameraEventNumberFromFilePath(line string) (int64, error) {
	re, err := regexp.Compile("/\\d+__")
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(strings.Trim(re.FindString(strings.Split(line, " saved to: ")[1]), "/_"), 10, 64)
}

func GetEventID(cameraNumber int64, cameraName string, cameraEventNumber int64) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fmt.Sprintf("%v-%v-%v", cameraNumber, cameraName, cameraEventNumber)))
}

func GetFilePath(line string) string {
	return strings.TrimSpace(strings.Split(line, " saved to: ")[1])
}

func ParseLine(timestamp time.Time, line string) (Event, error) {
	if len(strings.TrimSpace(line)) == 0 {
		return Event{}, fmt.Errorf("Empty line")
	}

	eventState := "unknown"

	if strings.Contains(line, "motion_detected: Motion detected - starting event") {
		eventState = "motion_start"
	} else if strings.Contains(line, "mlp_actions: End of event") {
		eventState = "motion_stop"
	} else if strings.Contains(line, "File of type 8 saved to") {
		eventState = "save_video"
	} else if strings.Contains(line, "File of type 1 saved to") {
		eventState = "save_image"
	}

	var cameraEventNumber int64

	var err error

	var filePath string

	if eventState == "motion_start" || eventState == "motion_stop" {
		cameraEventNumber, err = GetCameraEventNumberFromEndOfLine(line)
		if err != nil {
			return Event{}, err
		}
	} else if eventState == "save_video" || eventState == "save_image" {
		cameraEventNumber, err = GetCameraEventNumberFromFilePath(line)
		if err != nil {
			return Event{}, err
		}

		filePath = GetFilePath(line)
	} else {
		return Event{}, fmt.Errorf("unsupported line")
	}

	cameraNumber, err := GetCameraNumber(line)
	if err != nil {
		return Event{}, err
	}

	cameraName := GetCameraName(line)

	eventID := GetEventID(cameraNumber, cameraName, cameraEventNumber)

	if eventState == "unknown" {
		return Event{}, fmt.Errorf("unsupported event state")
	}

	return Event{
		Timestamp:         timestamp,
		EventID:           eventID,
		EventState:        eventState,
		CameraNumber:      cameraNumber,
		CameraName:        cameraName,
		CameraEventNumber: cameraEventNumber,
		FilePath:          filePath,
	}, nil
}
