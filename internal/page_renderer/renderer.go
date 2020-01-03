package page_renderer

import (
	"bytes"
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"text/template"
	"time"
)

type SummaryTableRowSeed struct {
	EventsURL  string
	EventsDate string
	EventCount string
}

func renderSummaryTableRows(eventsByDate map[time.Time][]event_store.Event) (string, error) {
	b := bytes.Buffer{}
	for eventsDate, events := range eventsByDate {
		t := template.New("SummaryTableRowSeed")

		t, err := t.Parse(SummaryTableRowHTML)
		if err != nil {
			return "", err
		}

		eventsSummaryTableRowSeed := SummaryTableRowSeed{
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

type SummarySeed struct {
	Title      string
	Now        string
	StyleSheet string
	TableRows  string
}

func RenderSummary(title string, eventsByDate map[time.Time][]event_store.Event, now time.Time) (string, error) {
	t := template.New("SummarySeed")

	t, err := t.Parse(SummaryHTML)
	if err != nil {
		return "", err
	}

	b := bytes.Buffer{}

	tableRows, err := renderSummaryTableRows(eventsByDate)
	if err != nil {
		return "", err
	}

	eventSummary := SummarySeed{
		Title:      title,
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

type PageTableRowSeed struct {
	EventID         string
	Timestamp       string
	Size            string
	CameraName      string
	HighResImageURL string
	LowResImageURL  string
	HighResVideoURL string
	LowResVideoURL  string
}

func renderPageTableRows(events []event_store.Event) (string, error) {
	b := bytes.Buffer{}
	for _, event := range events {
		t := template.New("PageTableRowSeed")

		t, err := t.Parse(PageTableRowHTML)
		if err != nil {
			return "", err
		}

		eventsTableRowSeed := PageTableRowSeed{
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

type PageSeed struct {
	Title      string
	EventsDate string
	Now        string
	StyleSheet string
	TableRows  string
}

func RenderPage(title string, events []event_store.Event, eventsDate, now time.Time) (string, error) {
	t := template.New("PageSeed")

	t, err := t.Parse(PageHTML)
	if err != nil {
		return "", err
	}

	b := bytes.Buffer{}

	tableRows, err := renderPageTableRows(events)
	if err != nil {
		return "", err
	}

	eventsSeed := PageSeed{
		Title:      title,
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
