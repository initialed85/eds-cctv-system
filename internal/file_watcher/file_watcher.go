package file_watcher

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/kylelemons/godebug/diff"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type FileWatcher struct {
	path     string
	callback func(time.Time, []string)
	watcher  *fsnotify.Watcher
	last     []byte
	stop     chan struct{}
}

func New(path string, callback func(time.Time, []string)) (FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return FileWatcher{}, err
	}

	last := make([]byte, 256)
	for i := 0; i < 256; i++ {
		last[i] = byte(i)
	}

	return FileWatcher{
		path:     path,
		callback: callback,
		watcher:  watcher,
		last:     last,
		stop:     make(chan struct{}),
	}, nil
}

func (f *FileWatcher) handle(timestamp time.Time) {
	this, err := ioutil.ReadFile(f.path)
	if err != nil {
		log.Printf("failed to read file because %v", err)

		return
	}

	if bytes.Compare(this, f.last) == 0 {
		log.Printf("file hasn't changed")

		return
	}

	lastLines := strings.Split(string(f.last), "\n")
	thisLines := strings.Split(string(this), "\n")

	allAdded := make([]string, 0)

	chunks := diff.DiffChunks(lastLines, thisLines)
	for _, chunk := range chunks {
		for _, added := range chunk.Added {
			allAdded = append(allAdded, added)
		}
	}

	log.Printf("calling callback w/ %v lines", len(allAdded))

	f.callback(timestamp, allAdded)

	f.last = this
}

func (f *FileWatcher) Watch() {
	defer func() {
		err := f.watcher.Close()
		if err != nil {
			log.Printf("closing watcher during defer caused %v", err)
		}
	}()

	fileDidNotPreviouslyExist := false

	for {
		_ = f.watcher.Remove(f.path)

		err := f.watcher.Add(f.path)
		if err != nil {
			if !fileDidNotPreviouslyExist {
				log.Printf("failed to add path to watcher because %v; will keep trying", err)
			}

			time.Sleep(time.Millisecond * 100)

			fileDidNotPreviouslyExist = true

			continue
		}

		info, err := os.Stat(f.path)
		if err != nil {
			log.Printf("failed to get stat for path; will try again later")

			continue
		}

		if info.IsDir() {
			log.Printf("path exists but is a directory; will try again later")

			continue
		}

		log.Printf("path exists; added file_watcher")

		if fileDidNotPreviouslyExist {
			log.Printf("path did not previously exist; calling handle")

			f.handle(time.Now())

			fileDidNotPreviouslyExist = false

			continue
		}

		restart := false

		for {
			if restart {
				break
			}

			select {
			case <-f.stop:
				log.Printf("stopping")

				return
			case event := <-f.watcher.Events:
				if event.Op != fsnotify.Write {
					log.Printf("non-write event caught; breaking")

					restart = true

					break
				}

				log.Printf("write event caught; calling handle")

				f.handle(time.Now())
			}
		}
	}
}

func (f *FileWatcher) Stop() {
	close(f.stop)
}
