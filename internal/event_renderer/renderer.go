package event_renderer

import (
	"bytes"
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"text/template"
	"time"
)

type EventsSummaryTableRowSeed struct {
	EventsURL  string
	EventsDate string
	EventCount string
}

func renderEventsSummaryTableRows(eventsByDate map[time.Time][]event_store.Event) (string, error) {
	b := bytes.Buffer{}
	for eventsDate, events := range eventsByDate {
		t := template.New("EventsSummaryTableRowSeed")

		t, err := t.Parse(EventsSummaryTableRowHTML)
		if err != nil {
			return "", err
		}

		eventsSummaryTableRowSeed := EventsSummaryTableRowSeed{
			EventsURL:  fmt.Sprintf("events_%v.html", eventsDate.Format("2006_01_02")),
			EventsDate: eventsDate.Format("2006-01-02"),
			EventCount: fmt.Sprintf("%v", len(events)),
		}

		err = t.Execute(&b, eventsSummaryTableRowSeed)
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

type EventsSummarySeed struct {
	Now        string
	StyleSheet string
	TableRows  string
}

func RenderEventsSummary(eventsByDate map[time.Time][]event_store.Event, now time.Time) (string, error) {
	t := template.New("EventsSummarySeed")

	t, err := t.Parse(EventsSummaryHTML)
	if err != nil {
		return "", err
	}

	b := bytes.Buffer{}

	tableRows, err := renderEventsSummaryTableRows(eventsByDate)
	if err != nil {
		return "", err
	}

	eventSummary := EventsSummarySeed{
		Now:        now.String(),
		StyleSheet: StyleSheet,
		TableRows:  tableRows,
	}

	err = t.Execute(&b, eventSummary)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

type EventsTableRowSeed struct {
	EventID         string
	Timestamp       string
	Size            string
	CameraName      string
	HighResImageURL string
	LowResImageURL  string
	HighResVideoURL string
	LowResVideoURL  string
}

func renderEventsTableRows(events []event_store.Event) (string, error) {
	b := bytes.Buffer{}
	for _, event := range events {
		t := template.New("EventsTableRowSeed")

		t, err := t.Parse(EventsTableRowHTML)
		if err != nil {
			return "", err
		}

		eventsTableRowSeed := EventsTableRowSeed{
			EventID:         event.EventID.String(),
			Timestamp:       event.Timestamp.String(),
			Size:            "?",
			CameraName:      event.CameraName,
			HighResImageURL: event.HighResImagePath,
			LowResImageURL:  event.LowResImagePath,
			HighResVideoURL: event.HighResVideoPath,
			LowResVideoURL:  event.LowResVideoPath,
		}

		err = t.Execute(&b, eventsTableRowSeed)
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

type EventsSeed struct {
	EventsDate string
	Now        string
	StyleSheet string
	TableRows  string
}

func RenderEvents(events []event_store.Event, eventsDate, now time.Time) (string, error) {
	t := template.New("EventsSeed")

	t, err := t.Parse(EventsHTML)
	if err != nil {
		return "", err
	}

	b := bytes.Buffer{}

	tableRows, err := renderEventsTableRows(events)
	if err != nil {
		return "", err
	}

	eventsSeed := EventsSeed{
		EventsDate: eventsDate.Format("2006-01-02"),
		Now:        now.String(),
		StyleSheet: StyleSheet,
		TableRows:  tableRows,
	}

	err = t.Execute(&b, eventsSeed)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
