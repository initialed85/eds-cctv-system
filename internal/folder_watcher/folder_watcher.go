package folder_watcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"strings"
	"time"
)

type FolderWatcher struct {
	path     string
	callback func(time.Time, string)
	watcher  *fsnotify.Watcher
	last     map[string]string
	stop     chan struct{}
}

func New(path string, callback func(time.Time, string)) (FolderWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return FolderWatcher{}, err
	}

	last := make([]byte, 256)
	for i := 0; i < 256; i++ {
		last[i] = byte(i)
	}

	return FolderWatcher{
		path:     path,
		callback: callback,
		watcher:  watcher,
		last:     make(map[string]string),
		stop:     make(chan struct{}),
	}, nil
}

func (f *FolderWatcher) handle(timestamp time.Time, path string) {
	parts := strings.Split(path, "_") // path = Segment_2019-12-15_22-04-28_Driveway.mp4

	suffix := parts[len(parts)-1]

	lastPath, ok := f.last[suffix]
	if !ok {
		f.last[suffix] = path

		return
	}

	if path != lastPath {
		log.Printf("calling callback w/ %v", path)

		f.callback(timestamp, lastPath)
	}

	f.last[suffix] = path
}

func (f *FolderWatcher) Watch() {
	defer func() {
		err := f.watcher.Close()
		if err != nil {
			log.Printf("closing watcher during defer caused %v", err)
		}
	}()

	folderDidNotPreviouslyExist := false

	for {
		_ = f.watcher.Remove(f.path)

		err := f.watcher.Add(f.path)
		if err != nil {
			if !folderDidNotPreviouslyExist {
				log.Printf("failed to add folder because %v; will keep trying", err)
			}

			time.Sleep(time.Millisecond * 100)

			folderDidNotPreviouslyExist = true

			continue
		}

		info, err := os.Stat(f.path)
		if err != nil {
			log.Printf("failed to get stat for path; will try again later")

			continue
		}

		if !info.IsDir() {
			log.Printf("path exists but is not directory; will try again later")

			continue
		}

		log.Printf("path exists; added watcher")

		folderDidNotPreviouslyExist = false

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
				if event.Op == fsnotify.Remove {
					log.Printf("remove event caught; breaking")

					restart = true

					break
				}

				if !strings.HasPrefix(event.Name, "Segment_") && !strings.HasSuffix(event.Name, ".mp4") {
					continue
				}

				if strings.Contains(event.Name, "-lowres") {
					continue
				}

				f.handle(time.Now(), event.Name)
			}
		}
	}
}

func (f *FolderWatcher) Stop() {
	close(f.stop)
}
