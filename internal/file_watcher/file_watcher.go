package file_watcher

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/kylelemons/godebug/diff"
	"io/ioutil"
	"strings"
	"time"
)

type FileWatcher struct {
	path     string
	callback func(string)
	watcher  *fsnotify.Watcher
	last     []byte
	stop     chan struct{}
}

func New(path string, callback func(string)) (FileWatcher, error) {
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

func (f *FileWatcher) handle() {
	this, err := ioutil.ReadFile(f.path)
	if err != nil {
		return
	}

	if bytes.Compare(this, f.last) == 0 {
		return
	}

	lastLines := strings.Split(string(f.last), "\n")
	thisLines := strings.Split(string(this), "\n")

	allAdded := ""

	chunks := diff.DiffChunks(lastLines, thisLines)
	for _, chunk := range chunks {
		for _, added := range chunk.Added {
			allAdded += added + "\n"
		}
	}

	if !strings.HasSuffix(string(this), "\n") {
		allAdded = strings.TrimRight(allAdded, "\n")
	}

	f.callback(allAdded)

	f.last = this
}

func (f *FileWatcher) Watch() {
	defer func() {
		err := f.watcher.Close()
		if err != nil {
		}
	}()

	fileDidNotPreviouslyExist := false

	for {
		err := f.watcher.Add(f.path)
		if err != nil {
			time.Sleep(time.Millisecond)

			fileDidNotPreviouslyExist = true

			continue
		}

		if fileDidNotPreviouslyExist {
			f.handle()

			continue
		}

		for {
			select {
			case <-f.stop:
				return
			case event := <-f.watcher.Events:
				if event.Op == fsnotify.Remove {
					break
				} else if event.Op != fsnotify.Write {
					continue
				}

				f.handle()
			}
		}
	}
}

func (f *FileWatcher) Stop() {
	close(f.stop)
}
