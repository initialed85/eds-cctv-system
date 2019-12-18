package event_store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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

func writeFile(data, path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSONLine(event Event, path string) error {
	data, err := MarshalJSONLine(event)
	if err != nil {
		return err
	}

	return writeFile(data, path)
}

func WriteJSONLines(events []Event, path string) error {
	output := ""

	for _, event := range events {
		data, err := MarshalJSONLine(event)
		if err != nil {
			return err
		}

		output += data
	}

	return writeFile(output, path)
}

func UnmarshalJSONLines(data string) ([]Event, error) {
	var events []Event

	data = fmt.Sprintf("[%v]", strings.Join(strings.Split(strings.TrimSpace(data), "\n"), ","))

	err := json.Unmarshal([]byte(data), &events)
	if err != nil {
		return []Event{}, err
	}

	return events, nil
}

func ReadJSONLines(path string) ([]Event, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []Event{}, err
	}

	events, err := UnmarshalJSONLines(string(data))
	if err != nil {
		return []Event{}, err
	}

	return events, err
}
