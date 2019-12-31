package file_write_folder_watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/initialed85/eds-cctv-system/internal/file_watcher"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type Watcher struct {
	path        string
	callback    func(time.Time, string)
	fileWatcher file_watcher.Watcher
	last        map[string]string
}

func New(path string, callback func(time.Time, string)) (Watcher, error) {
	last := make([]byte, 256)
	for i := 0; i < 256; i++ {
		last[i] = byte(i)
	}

	w := Watcher{
		path:     path,
		callback: callback,
		last:     make(map[string]string),
	}

	fileWatcher, err := file_watcher.New(path, w.handle)
	if err != nil {
		return Watcher{}, err
	}

	w.fileWatcher = fileWatcher

	return w, nil
}

func (w *Watcher) handle(event fsnotify.Event, timestamp time.Time) bool {
	if event.Op != fsnotify.Create {
		return true
	}

	_, fileName := filepath.Split(event.Name)

	if !strings.HasPrefix(fileName, "Segment_") {
		log.Printf("%v doesn't have Segment_ prefix", fileName)

		return true
	}

	if !strings.HasSuffix(fileName, ".mp4") {
		log.Printf("%v doesn't have .mp4 suffix", fileName)

		return true
	}

	if strings.Contains(fileName, "-lowres") {
		log.Printf("%v contains -lowres", fileName)

		return true
	}

	path := event.Name

	parts := strings.Split(path, "_") // path = Segment_2019-12-15_22-04-28_Driveway.mp4

	suffix := parts[len(parts)-1]

	lastPath, ok := w.last[suffix]
	if !ok {
		w.last[suffix] = path

		return true
	}

	pathChanged := path != lastPath

	log.Printf("handled path=%v, suffix=%v, lastPath=%v, pathChanged=%v", path, suffix, lastPath, pathChanged)

	if pathChanged {
		log.Printf("calling callback w/ %v", lastPath)

		w.callback(timestamp, lastPath)
	}

	w.last[suffix] = path

	return true
}

func (w *Watcher) Watch() {
	w.fileWatcher.Watch()
}

func (w *Watcher) Stop() {
	w.fileWatcher.Stop()
}
