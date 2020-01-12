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

func baseWriteFile(data, path string, flag int) error {
	f, err := os.OpenFile(path, flag, 0644)
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

func writeFile(data, path string) error {
	return baseWriteFile(data, path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY)
}

func appendFile(data, path string) error {
	return baseWriteFile(data, path, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
}

func buildJSONLines(events []Event) (string, error) {
	output := ""

	for _, event := range events {
		data, err := MarshalJSONLine(event)
		if err != nil {
			return "", err
		}

		output += data
	}

	return output, nil
}

func WriteJSONLines(events []Event, path string) error {
	output, err := buildJSONLines(events)
	if err != nil {
		return err
	}

	return writeFile(output, path)
}

func AppendJSONLines(events []Event, path string) error {
	output, err := buildJSONLines(events)
	if err != nil {
		return err
	}

	return appendFile(output, path)
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
