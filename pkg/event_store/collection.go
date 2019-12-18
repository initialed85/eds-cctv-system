package event_store

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Event struct {
	EventID          uuid.UUID `json:"event_id"`
	Timestamp        time.Time `json:"timestamp"`
	HighResImagePath string    `json:"high_res_image_path"`
	LowResImagePath  string    `json:"low_res_image_path"`
	HighResVideoPath string    `json:"high_res_video_path"`
	LowResVideoPath  string    `json:"low_res_video_path"`
}

func NewEvent(timestamp time.Time, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath string) Event {
	return Event{
		EventID: uuid.NewSHA1(
			uuid.NameSpaceDNS,
			[]byte(fmt.Sprintf("%v, %v, %v, %v, %v", timestamp, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath)),
		),
		Timestamp:        timestamp,
		HighResImagePath: highResImagePath,
		LowResImagePath:  lowResImagePath,
		HighResVideoPath: highResVideoPath,
		LowResVideoPath:  lowResImagePath,
	}
}

type EventCollection struct {
	path   string
	events map[uuid.UUID]Event
	mu     sync.Mutex
}

func NewEventCollection(path string) EventCollection {
	e := EventCollection{
		path:   path,
		events: make(map[uuid.UUID]Event),
	}

	_ = e.Read()

	return e
}

func (e *EventCollection) Read() error {
	events, err := ReadJSONLines(e.path)
	if err != nil {
		return err
	}

	e.mu.Lock()
	for _, event := range events {
		e.events[event.EventID] = event
	}
	e.mu.Unlock()

	return nil
}

func (e *EventCollection) Write() {
	events := make([]Event, 0)

	e.mu.Lock()
	for _, event := range e.events {
		events = append(events, event)
	}
	e.mu.Unlock()
}

func (e *EventCollection) Add(event Event) {
	e.mu.Lock()
	e.events[event.EventID] = event
	e.mu.Unlock()
}

func (e *EventCollection) Remove(eventID uuid.UUID) error {
	e.mu.Lock()
	_, ok := e.events[eventID]
	e.mu.Unlock()

	if !ok {
		return fmt.Errorf("event with EventID=%v did not exist", eventID)
	}

	return nil
}

func (e *EventCollection) Get(eventID uuid.UUID) (Event, error) {
	e.mu.Lock()
	event, ok := e.events[eventID]
	e.mu.Unlock()

	if !ok {
		return Event{}, fmt.Errorf("event with EventID=%v did not exist", eventID)
	}

	return event, nil
}

func (e *EventCollection) Pop(eventID uuid.UUID) (Event, error) {
	e.mu.Lock()
	event, ok := e.events[eventID]
	if ok {
		delete(e.events, eventID)
	}
	e.mu.Unlock()

	if !ok {
		return Event{}, fmt.Errorf("event with EventID=%v did not exist", eventID)
	}

	return event, nil
}

func (e *EventCollection) Len() int {
	e.mu.Lock()
	length := len(e.events)
	e.mu.Unlock()

	return length
}
