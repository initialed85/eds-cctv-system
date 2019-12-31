package event_store

import (
	"fmt"
	"github.com/google/uuid"
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

	_ = e.Read()

	return e
}

func (s *Store) GetPath() string {
	return s.path
}

func (s *Store) Read() error {
	events, err := ReadJSONLines(s.path)
	if err != nil {
		return err
	}

	s.mu.Lock()
	for _, event := range events {
		s.events[event.EventID] = event
	}
	s.mu.Unlock()

	return nil
}

func (s *Store) Write() error {
	events := make([]Event, 0)

	s.mu.Lock()
	for _, event := range s.events {
		events = append(events, event)
	}
	s.mu.Unlock()

	err := WriteJSONLines(events, s.path)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Len() int {
	s.mu.Lock()
	length := len(s.events)
	s.mu.Unlock()

	return length
}

func (s *Store) Add(event Event) {
	s.mu.Lock()
	s.events[event.EventID] = event
	s.mu.Unlock()
}

func (s *Store) GetAll() []Event {
	if s.Len() == 0 {
		return []Event{}
	}

	var events []Event

	s.mu.Lock()
	for _, event := range s.events {
		events = append(events, event)
	}
	s.mu.Unlock()

	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Timestamp.Unix() < events[j].Timestamp.Unix()
	})

	return events
}

func (s *Store) GetAllByDate() map[time.Time][]Event {
	allEvents := s.GetAll()

	eventsByDate := make(map[time.Time][]Event)

	for _, event := range allEvents {
		date := event.Timestamp.Truncate(time.Hour * 24)

		_, ok := eventsByDate[date]
		if !ok {
			eventsByDate[date] = make([]Event, 0)
		}

		eventsByDate[date] = append(eventsByDate[date], event)
	}

	for date := range eventsByDate {
		sort.SliceStable(eventsByDate[date], func(i, j int) bool {
			return eventsByDate[date][i].Timestamp.Unix() < eventsByDate[date][j].Timestamp.Unix()
		})
	}

	return eventsByDate
}

func (s *Store) Get(eventID uuid.UUID) (Event, error) {
	s.mu.Lock()
	event, ok := s.events[eventID]
	s.mu.Unlock()

	if !ok {
		return Event{}, fmt.Errorf("event with EventID=%v did not exist", eventID)
	}

	return event, nil
}

func (s *Store) Pop(eventID uuid.UUID) (Event, error) {
	s.mu.Lock()
	event, ok := s.events[eventID]
	if ok {
		delete(s.events, eventID)
	}
	s.mu.Unlock()

	if !ok {
		return Event{}, fmt.Errorf("event with EventID=%v did not exist", eventID)
	}

	return event, nil
}

func (s *Store) Remove(eventID uuid.UUID) error {
	_, err := s.Pop(eventID)

	return err
}