package file_diff_file_watcher

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/initialed85/eds-cctv-system/internal/file_watcher"
	"github.com/kylelemons/godebug/diff"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type Watcher struct {
	path        string
	callback    func(time.Time, []string)
	fileWatcher file_watcher.Watcher
	last        []byte
}

func New(path string, callback func(time.Time, []string)) (Watcher, error) {
	last := make([]byte, 256)
	for i := 0; i < 256; i++ {
		last[i] = byte(i)
	}

	w := Watcher{
		path:     path,
		callback: callback,
		last:     last,
	}

	dir, _ := filepath.Split(path)

	fileWatcher, err := file_watcher.New(dir, w.handle)
	if err != nil {
		return Watcher{}, err
	}

	w.fileWatcher = fileWatcher

	return w, nil
}

func (w *Watcher) handle(event fsnotify.Event, timestamp time.Time) bool {
	if event.Name != w.path {
		return true
	}

	if event.Op == fsnotify.Create {
		w.callback(timestamp, make([]string, 0))

		return true
	}

	if event.Op != fsnotify.Write {
		log.Printf("non-write event caught; breaking")

		return false
	}

	this, err := ioutil.ReadFile(event.Name)
	if err != nil {
		log.Printf("failed to read file because %v", err)

		return true
	}

	if bytes.Compare(this, w.last) == 0 {
		log.Printf("file hasn't changed")

		return true
	}

	lastLines := strings.Split(string(w.last), "\n")
	thisLines := strings.Split(string(this), "\n")

	allAdded := make([]string, 0)

	chunks := diff.DiffChunks(lastLines, thisLines)
	for _, chunk := range chunks {
		for _, added := range chunk.Added {
			allAdded = append(allAdded, added)
		}
	}

	log.Printf("calling callback w/ %v lines", len(allAdded))

	go w.callback(timestamp, allAdded)

	w.last = this

	return true
}

func (w *Watcher) Watch() {
	w.fileWatcher.Watch()
}

func (w *Watcher) Stop() {
	w.fileWatcher.Stop()
}
