package file_watcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"time"
)

type Watcher struct {
	path     string
	callback func(fsnotify.Event, time.Time) bool
	watcher  *fsnotify.Watcher
	stopped  bool
	stop     chan struct{}
}

func New(path string, callback func(fsnotify.Event, time.Time) bool) (Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return Watcher{}, err
	}

	last := make([]byte, 256)
	for i := 0; i < 256; i++ {
		last[i] = byte(i)
	}

	return Watcher{
		path:     path,
		callback: callback,
		watcher:  watcher,
		stopped:  false,
		stop:     make(chan struct{}),
	}, nil
}

func (w *Watcher) Watch() {
	defer func() {
		err := w.watcher.Close()
		if err != nil {
			log.Printf("closing watcher during defer caused %v", err)
		}
	}()

	pathDidNotPreviouslyExist := false

	for {
		if w.stopped {
			log.Printf("stopping")

			return
		}

		_ = w.watcher.Remove(w.path)

		err := w.watcher.Add(w.path)
		if err != nil {
			if !pathDidNotPreviouslyExist {
				log.Printf("failed to add folder because %v; will keep trying", err)
			}

			time.Sleep(time.Millisecond * 100)

			pathDidNotPreviouslyExist = true

			continue
		}

		info, err := os.Stat(w.path)
		if err != nil {
			log.Printf("failed to get stat for path; will try again later")

			continue
		}

		if !info.IsDir() {
			log.Printf("path exists but is not directory; will try again later")

			continue
		}

		log.Printf("path exists; added watcher")

		pathDidNotPreviouslyExist = false

		restart := false

		for {
			if restart {
				break
			}

			select {
			case <-w.stop:
				log.Printf("stopping")

				return
			case event := <-w.watcher.Events:
				if event.Name == w.path && event.Op == fsnotify.Remove {
					log.Printf("remove event caught; breaking")

					restart = true

					break
				}

				restart = !w.callback(event, time.Now())
			}
		}
	}
}

func (w *Watcher) Stop() {
	w.stopped = true

	close(w.stop)
}
