package motion_log_event_streamer

import (
	"eds-cctv-system/internal/event_streamer"
	"eds-cctv-system/internal/file_watcher"
	"eds-cctv-system/internal/motion_log_parser"
	"encoding/json"
	"log"
	"time"
)

type MotionLogEventStreamer struct {
	fileWatcher   file_watcher.FileWatcher
	eventStreamer *event_streamer.EventStreamer
}

func New(filePath string, port int64) (MotionLogEventStreamer, error) {
	motionLogEventStreamer := MotionLogEventStreamer{}

	fileWatcher, err := file_watcher.New(filePath, motionLogEventStreamer.fileWatchCallback)
	if err != nil {
		return motionLogEventStreamer, err
	}

	eventStreamer := event_streamer.New(port)

	motionLogEventStreamer.fileWatcher = fileWatcher
	motionLogEventStreamer.eventStreamer = eventStreamer

	return motionLogEventStreamer, nil
}

func (m *MotionLogEventStreamer) fileWatchCallback(timestamp time.Time, lines []string) {
	for _, line := range lines {
		event, err := motion_log_parser.ParseLine(timestamp, line)
		if err != nil {
			log.Printf("Error %v while parsing '%v'; skipping...", err, line)

			continue
		}

		jsonEvent, err := json.Marshal(event)
		if err != nil {
			log.Printf("Error %v while marshalling %v; skipping...", err, event)

			continue
		}

		m.eventStreamer.AddMessage(string(jsonEvent))
	}
}

func (m *MotionLogEventStreamer) Start() {
	go func() {
		err := m.eventStreamer.Listen()
		if err != nil {
			log.Fatalf("Error %v while trying call %v", err, m.eventStreamer.Listen())
		}
	}()

	go m.fileWatcher.Watch()

	time.Sleep(time.Second)
}

func (m *MotionLogEventStreamer) Stop() error {
	err := m.eventStreamer.Stop()

	m.fileWatcher.Stop()

	return err
}
