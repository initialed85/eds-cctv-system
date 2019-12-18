package event_persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Event struct {
	Timestamp        time.Time `json:"timestamp"`
	HighResImagePath string    `json:"high_res_image_path"`
	LowResImagePath  string    `json:"low_res_image_path"`
	HighResVideoPath string    `json:"high_res_video_path"`
	LowResVideoPath  string    `json:"low_res_video_path"`
}

func MarshalJSONLine(event Event) (string, error) {
	var events []Event

	events = append(events, event)

	b, err := json.Marshal(events)
	if err != nil {
		return "", err
	}

	data := string(b)

	return fmt.Sprintf("%v\n", data[1:len(data)-1]), nil
}

func WriteJSONLine(event Event, path string) error {
	jsonLine, err := MarshalJSONLine(event)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = f.WriteString(jsonLine)
	if err != nil {
		return err
	}

	return nil
}
