package event_store_updater

import (
	"github.com/fsnotify/fsnotify"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"github.com/initialed85/eds-cctv-system/internal/file_watcher"
	"log"
	"path/filepath"
	"time"
)

type Updater struct {
	store    *event_store.Store
	callback func(*event_store.Store)
	watcher  file_watcher.Watcher
}

func New(store *event_store.Store, callback func(*event_store.Store)) (Updater, error) {
	u := Updater{
		store:    store,
		callback: callback,
	}

	dir, _ := filepath.Split(store.GetPath())

	watcher, err := file_watcher.New(dir, u.handle)
	if err != nil {
		return Updater{}, err
	}

	u.watcher = watcher

	return u, nil
}

func (u *Updater) handle(event fsnotify.Event, timestamp time.Time) bool {
	if event.Name != u.store.GetPath() {
		return true
	}

	err := u.store.Read()
	if err != nil {
		log.Printf("failed to read store because: %v", err)
	}

	u.callback(u.store)

	return true
}

func (u *Updater) Watch() {
	u.watcher.Watch()
}

func (u *Updater) Stop() {
	u.watcher.Stop()
}
