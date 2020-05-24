package event_store

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"sort"
	"sync"
	"time"
)

type Store struct {
	path   string
	events map[uuid.UUID]Event
	mu     sync.Mutex
}

func NewStore(path string) Store {
	e := Store{
		path:   path,
		events: make(map[uuid.UUID]Event),
	}

	return e
}

func (s *Store) GetPath() string {
	return s.path
}

func (s *Store) Read() error {
	log.Printf("reading from %v", s.path)

	events, err := ReadJSONLines(s.path)
	if err != nil {
		return err
	}

	log.Printf("read %v JSON Lines from %v", len(events), s.path)

	s.mu.Lock()
	for _, event := range events {
		s.events[event.EventID] = event
	}
	s.mu.Unlock()

	log.Printf("read %v unique events from JSON Lines", s.Len())

	return nil
}

func (s *Store) Write() error {
	events := s.getAll()

	log.Printf("writing %v events to %v", len(events), s.path)

	return WriteJSONLines(events, s.path)
}

func (s *Store) Append() error {
	events := s.getAll()

	log.Printf("appending %v events to %v", len(events), s.path)

	return AppendJSONLines(events, s.path)
}

func (s *Store) Len() int {
	s.mu.Lock()
	length := len(s.events)
	s.mu.Unlock()

	return length
}

func (s *Store) Overwrite(event Event) {
	log.Printf("adding/overwriting %+v", event)

	s.mu.Lock()
	s.events[event.EventID] = event
	s.mu.Unlock()
}

func (s *Store) Add(event Event) {
	pathComparison := fmt.Sprintf(
		"%v-%v-%v-%v",
		event.HighResImagePath,
		event.HighResVideoPath,
		event.LowResImagePath,
		event.LowResVideoPath,
	)

	for _, otherEvent := range s.getAll() {
		otherPathComparison := fmt.Sprintf(
			"%v-%v-%v-%v",
			otherEvent.HighResImagePath,
			otherEvent.HighResVideoPath,
			otherEvent.LowResImagePath,
			otherEvent.LowResVideoPath,
		)

		if pathComparison == otherPathComparison {
			return
		}
	}

	s.Overwrite(event)
}

func (s *Store) getAll() []Event {
	if s.Len() == 0 {
		return []Event{}
	}

	var events []Event

	s.mu.Lock()
	for _, event := range s.events {
		events = append(events, event)
	}
	s.mu.Unlock()

	return events
}

func (s *Store) GetAll() []Event {
	events := s.getAll()

	log.Printf("got %v events", len(events))

	return events
}

func (s *Store) GetAllDescending() []Event {
	events := s.getAll()

	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Timestamp.Unix() > events[j].Timestamp.Unix()
	})

	log.Printf("got %v events", len(events))

	return events
}

func (s *Store) GetAllDescendingByDateDescending(cameraName string) map[time.Time][]Event {
	allEvents := s.GetAllDescending()

	eventsByDate := make(map[time.Time][]Event)

	for _, event := range allEvents {
		date, _ := time.Parse("2006-01-02", event.Timestamp.Format("2006-01-02"))

		_, ok := eventsByDate[date]
		if !ok {
			eventsByDate[date] = make([]Event, 0)
		}

		eventsByDate[date] = append(eventsByDate[date], event)
	}

	for date := range eventsByDate {
		sort.SliceStable(eventsByDate[date], func(i, j int) bool {
			return eventsByDate[date][i].Timestamp.Unix() > eventsByDate[date][j].Timestamp.Unix()
		})
	}

	keys := make([]time.Time, 0)
	for key := range eventsByDate {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].Unix() > keys[j].Unix()
	})

	sortedEventsByDate := make(map[time.Time][]Event)
	for _, key := range keys {
		sortedEventsByDate[key] = eventsByDate[key]
	}

	log.Printf("got %v events across %v dates", len(allEvents), len(sortedEventsByDate))

	return sortedEventsByDate
}

func (s *Store) GetByUUID(eventID uuid.UUID) (Event, error) {
	s.mu.Lock()
	event, ok := s.events[eventID]
	s.mu.Unlock()

	if !ok {
		return Event{}, fmt.Errorf("event with EventID=%v did not exist", eventID)
	}

	log.Printf("got %+v", event)

	return event, nil
}
