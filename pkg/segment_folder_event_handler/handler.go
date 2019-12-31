package segment_folder_event_handler

import (
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/common"
	"github.com/initialed85/eds-cctv-system/internal/file_converter"
	"github.com/initialed85/eds-cctv-system/internal/file_write_folder_watcher"
	"github.com/initialed85/eds-cctv-system/internal/thumbnail_creator"
	"log"
	"path/filepath"
	"strings"
	"time"
)

func getImagePath(path string) string {
	extension := filepath.Ext(path)

	parts := strings.Split(path, extension)

	part := strings.Join(parts[0:len(parts)-1], extension)

	return fmt.Sprintf("%v%v", part, ".jpg")
}

type Handler struct {
	folderWatcher *file_write_folder_watcher.Watcher
	callback      func(time.Time, string, string, string, string) error
}

func New(folderPath string, callback func(time.Time, string, string, string, string) error) (Handler, error) {
	s := Handler{}

	folderWatcher, err := file_write_folder_watcher.New(folderPath, s.folderWatcherCallback)
	if err != nil {
		return Handler{}, err
	}

	s.folderWatcher = &folderWatcher
	s.callback = callback

	return s, nil
}

func (h *Handler) folderWatcherCallback(timestamp time.Time, highResVideoPath string) {
	highResImagePath := getImagePath(highResVideoPath)

	err := thumbnail_creator.CreateThumbnail(highResVideoPath, highResImagePath)
	if err != nil {
		log.Printf("failed to create thumbnail for %v because: %v", highResImagePath, err)

		return
	}

	log.Printf("created thumbnail for %v at %v", highResVideoPath, highResImagePath)

	lowResVideoPath := common.GetLowResPath(highResVideoPath)
	_, stderr, err := file_converter.ConvertVideo(highResVideoPath, lowResVideoPath, 640, 480)
	if err != nil {
		log.Printf("failed to convert %v to %v because %v; stderr=%v", highResVideoPath, lowResVideoPath, err, stderr)

		return
	}

	log.Printf("converted %v to %v", highResVideoPath, lowResVideoPath)

	lowResImagePath := common.GetLowResPath(highResImagePath)
	_, stderr, err = file_converter.ConvertImage(highResImagePath, lowResImagePath, 640, 480)
	if err != nil {
		log.Printf("failed to convert %v to %v because %v; stderr=%v", highResImagePath, lowResImagePath, err, stderr)

		return
	}

	log.Printf("converted %v to %v", highResImagePath, lowResImagePath)

	err = h.callback(timestamp, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath)
	if err != nil {
		log.Printf("failed to call callback with timestamp=%v, highResImagePath=%v, lowResImagePath=%v, highResVideoPath=%v, lowResVideoPath=%v because %v", timestamp, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath, err)

		return
	}

	log.Printf("called callback with highResImagePath=%v, lowResImagePath=%v, highResVideoPath=%v, lowResVideoPath=%v", highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath)
}

func (h *Handler) Start() {
	go h.folderWatcher.Watch()

	time.Sleep(time.Second)
}

func (h *Handler) Stop() error {
	h.folderWatcher.Stop()

	return nil
}
