package event_store_updater_page_renderer

import (
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"github.com/initialed85/eds-cctv-system/internal/event_store_updater"
	"github.com/initialed85/eds-cctv-system/internal/page_renderer"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	fileNamePrefix = "events"
	fileNameSuffix = "html"
)

func cleanFolder(path string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		_, file := filepath.Split(path)

		if !strings.HasSuffix(file, fileNameSuffix) {
			return nil
		}

		if !(strings.HasPrefix(file, fileNamePrefix)) {
			return nil
		}

		err = os.Remove(path)
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk(path, walkFn)
	if err != nil {
		return err
	}

	return nil
}

func writeFile(path, data string) error {
	return ioutil.WriteFile(path, []byte(data), 0644)
}

func truncatePath(path string) string {
	_, file := filepath.Split(path)

	return file
}

func truncatePaths(events []event_store.Event) []event_store.Event {
	newEvents := make([]event_store.Event, 0)

	for _, event := range events {
		newEvents = append(
			newEvents,
			event_store.Event{
				EventID:          event.EventID,
				Timestamp:        event.Timestamp,
				CameraName:       event.CameraName,
				HighResImagePath: truncatePath(event.HighResImagePath),
				LowResImagePath:  truncatePath(event.LowResImagePath),
				HighResVideoPath: truncatePath(event.HighResVideoPath),
				LowResVideoPath:  truncatePath(event.LowResVideoPath),
			},
		)
	}

	return newEvents
}

type Renderer struct {
	summaryTitle, title string
	store               event_store.Store
	updater             event_store_updater.Updater
	renderPath          string
}

func New(summaryTitle, title, storePath, renderPath string) (Renderer, error) {
	r := Renderer{
		summaryTitle: summaryTitle,
		title:        title,
		store:        event_store.NewStore(storePath),
		renderPath:   renderPath,
	}

	err := r.store.Read()
	if err != nil {
		return Renderer{}, err
	}

	updater, err := event_store_updater.New(r.store, r.callback)
	if err != nil {
		return Renderer{}, err
	}

	r.updater = updater

	return r, nil
}

func (r *Renderer) callback(store event_store.Store) {
	err := cleanFolder(r.renderPath)
	if err != nil {
		log.Printf("failed to call cleanFolder because: %v", err)

		return
	}

	now := time.Now()

	eventsByDate := store.GetAllByDate()

	for eventsDate, events := range eventsByDate {
		eventsHTML, err := page_renderer.RenderPage(r.title, truncatePaths(events), eventsDate, now)
		if err != nil {
			log.Printf("failed to call RenderPage for %v because: %v", eventsDate, err)

			continue
		}

		path := filepath.Join(
			r.renderPath,
			fmt.Sprintf(
				"%v_%v.%v",
				fileNamePrefix,
				eventsDate.Format("2006_01_02"),
				fileNameSuffix,
			),
		)

		err = writeFile(path, eventsHTML)
		if err != nil {
			log.Printf("failed to call writeFile for %v because: %v", path, err)

			continue
		}
	}

	eventsSummaryHTML, err := page_renderer.RenderSummary(r.summaryTitle, eventsByDate, now)
	if err != nil {
		log.Printf("failed to call RenderSummary because: %v", err)

		return
	}

	path := filepath.Join(
		r.renderPath,
		fmt.Sprintf("%v.%v", fileNamePrefix, fileNameSuffix),
	)

	err = writeFile(path, eventsSummaryHTML)
	if err != nil {
		log.Printf("failed to call writeFile for %v because: %v", path, err)

		return
	}
}

func (r *Renderer) Start() {
	go r.updater.Watch()

	// run once to pre-populate
	r.callback(r.store)
}

func (r *Renderer) Stop() {
	r.updater.Stop()
}
