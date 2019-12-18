package folder_watcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
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

	pathChanged := path != lastPath

	log.Printf("handled path=%v, suffix=%v, lastPath=%v, pathChanged=%v", path, suffix, lastPath, pathChanged)

	if pathChanged {
		log.Printf("calling callback w/ %v", lastPath)

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

				if event.Op != fsnotify.Create {
					continue
				}

				_, fileName := filepath.Split(event.Name)

				if !strings.HasPrefix(fileName, "Segment_") {
					log.Printf("%v doesn't have Segment_ prefix", fileName)

					continue
				}

				if !strings.HasSuffix(fileName, ".mp4") {
					log.Printf("%v doesn't have .mp4 suffix", fileName)

					continue
				}

				if strings.Contains(fileName, "-lowres") {
					log.Printf("%v contains -lowres", fileName)

					continue
				}

				log.Printf("calling handle with path=%v", event.Name)

				f.handle(time.Now(), event.Name)
			}
		}
	}
}

func (f *FolderWatcher) Stop() {
	close(f.stop)
}
