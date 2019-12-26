package event_renderer

import (
	"bytes"
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"text/template"
	"time"
)

type EventsSummaryTableRow struct {
	EventsURL  string
	EventsDate string
	EventCount string
}

func renderEventsSummaryTableRows(eventsByDate map[time.Time][]event_store.Event) (string, error) {
	b := bytes.Buffer{}
	for date, events := range eventsByDate {
		t := template.New("EventsSummaryTableRow")

		t, err := t.Parse(EventsSummaryTableRowHTML)
		if err != nil {
			return "", err
		}

		eventsSummaryTableRow := EventsSummaryTableRow{
			EventsURL:  fmt.Sprintf("events_%v.html", date.Format("2006_01_02")),
			EventsDate: date.Format("2006-01-02"),
			EventCount: fmt.Sprintf("%v", len(events)),
		}

		err = t.Execute(&b, eventsSummaryTableRow)
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

type EventsSummary struct {
	Now        string
	StyleSheet string
	TableRows  string
}

func RenderEventsSummary(eventsByDate map[time.Time][]event_store.Event, now time.Time) (string, error) {
	t := template.New("EventsSummary")

	t, err := t.Parse(EventsSummaryHTML)
	if err != nil {
		return "", err
	}

	b := bytes.Buffer{}

	tableRows, err := renderEventsSummaryTableRows(eventsByDate)
	if err != nil {
		return "", err
	}

	eventSummary := EventsSummary{
		StyleSheet: StyleSheet,
		Now:        now.String(),
		TableRows:  tableRows,
	}

	err = t.Execute(&b, eventSummary)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

type EventsTableRow struct {
	EventID         string
	CameraID        string
	Timestamp       string
	Size            string
	Camera          string
	HighResImageURL string
	LowResImageURL  string
	HighResVideoURL string
	LowResVideoURL  string
}

type Events struct {
	EventsDate string
	Now        string
	StyleSheet string
	TableRows  string
}

func RenderEvents(events []event_store.Event, eventsDate, now time.Time) (string, error) {

}