package motion_log_event_handler

import (
	"github.com/initialed85/eds-cctv-system/internal/common"
	"github.com/initialed85/eds-cctv-system/internal/file_converter"
	"github.com/initialed85/eds-cctv-system/internal/file_watcher"
	"github.com/initialed85/eds-cctv-system/internal/motion_log"
	"github.com/google/uuid"
	"log"
	"time"
)

type RelatedEvent struct {
	MotionStart     *motion_log.Event
	SaveVideo       *motion_log.Event
	SaveImage       *motion_log.Event
	MotionStop      *motion_log.Event
	LowResVideoPath string
	LowResImagePath string
}

func (r *RelatedEvent) IsComplete() bool {
	return r.MotionStart != nil && r.SaveVideo != nil && r.SaveImage != nil && r.MotionStop != nil
}

type MotionLogEventHandler struct {
	fileWatcher   file_watcher.FileWatcher
	relatedEvents map[uuid.UUID]RelatedEvent
}

func New(filePath string) (MotionLogEventHandler, error) {
	m := MotionLogEventHandler{
		relatedEvents: make(map[uuid.UUID]RelatedEvent),
	}

	fileWatcher, err := file_watcher.New(filePath, m.fileWatcherCallback)
	if err != nil {
		return MotionLogEventHandler{}, err
	}

	m.fileWatcher = fileWatcher

	return m, nil
}

func (m *MotionLogEventHandler) fileWatcherCallback(timestamp time.Time, lines []string) {
	for _, line := range lines {
		event, err := motion_log.ParseLine(timestamp, line)
		if err != nil {
			log.Printf("Error %v while parsing '%v'; skipping...", err, line)

			continue
		}

		relatedEvent, ok := m.relatedEvents[event.EventId]
		if !ok {
			relatedEvent = RelatedEvent{}
		}

		if event.EventState == motion_log.MotionStart {
			relatedEvent.MotionStart = &event
		} else if event.EventState == motion_log.SaveVideo {
			relatedEvent.SaveVideo = &event
		} else if event.EventState == motion_log.SaveImage {
			relatedEvent.SaveImage = &event
		} else if event.EventState == motion_log.MotionStop {
			relatedEvent.MotionStop = &event
		}

		m.relatedEvents[event.EventId] = relatedEvent

		if relatedEvent.IsComplete() {
			highResVideoPath := relatedEvent.SaveVideo.FilePath
			lowResVideoPath := common.GetLowResPath(highResVideoPath)
			_, stderr, err := file_converter.ConvertVideo(highResVideoPath, lowResVideoPath, 640, 480)
			if err != nil {
				log.Printf("failed to convert %v to %v because %v; stderr=%v", highResVideoPath, lowResVideoPath, err, stderr)
			}

			log.Printf("converted %v to %v", highResVideoPath, lowResVideoPath)

			highResImagePath := relatedEvent.SaveImage.FilePath
			lowResImagePath := common.GetLowResPath(highResImagePath)
			_, stderr, err = file_converter.ConvertImage(highResImagePath, lowResImagePath, 640, 480)
			if err != nil {
				log.Printf("failed to convert %v to %v because %v; stderr=%v", highResImagePath, lowResImagePath, err, stderr)
			}

			log.Printf("converted %v to %v", highResImagePath, lowResImagePath)

			delete(m.relatedEvents, event.EventId)
		}
	}
}

func (m *MotionLogEventHandler) Start() {
	go m.fileWatcher.Watch()

	time.Sleep(time.Second)
}

func (m *MotionLogEventHandler) Stop() error {
	m.fileWatcher.Stop()

	return nil
}
