package event_renderer

import (
	"bytes"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"html/template"
	"time"
)

type EventsSummaryTableRow struct {
}

func renderEventsSummaryTableRows(eventsByDate map[time.Time][]event_store.Event) (string, error) {
	b := bytes.Buffer{}
	for date, event := range eventsByDate {
		_ = date
		_ = event

		t := template.New("EventsSummaryTableRow")

		t, err := t.Parse(EventsSummaryTableRowHTML)
		if err != nil {
			return "", err
		}

		eventsSummaryTableRow := EventsSummaryTableRow{}

		err = t.Execute(&b, eventsSummaryTableRow)
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

type EventsSummary struct {
	StyleSheet string
	Now        string
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
