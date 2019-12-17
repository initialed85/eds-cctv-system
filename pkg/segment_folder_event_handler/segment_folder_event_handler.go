package segment_folder_event_handler

import (
	"github.com/initialed85/eds-cctv-system/internal/common"
	"github.com/initialed85/eds-cctv-system/internal/file_converter"
	"github.com/initialed85/eds-cctv-system/internal/folder_watcher"
	"github.com/initialed85/eds-cctv-system/internal/thumbnail_creator"
	"fmt"
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

type SegmentFolderEventHandler struct {
	folderWatcher folder_watcher.FolderWatcher
}

func New(folderPath string) (SegmentFolderEventHandler, error) {
	m := SegmentFolderEventHandler{}

	folderWatcher, err := folder_watcher.New(folderPath, m.folderWatcherCallback)
	if err != nil {
		return SegmentFolderEventHandler{}, err
	}

	m.folderWatcher = folderWatcher

	return m, nil
}

func (s *SegmentFolderEventHandler) folderWatcherCallback(timestamp time.Time, highResVideoPath string) {
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
	}

	log.Printf("converted %v to %v", highResVideoPath, lowResVideoPath)

	lowResImagePath := common.GetLowResPath(highResImagePath)
	_, stderr, err = file_converter.ConvertImage(highResImagePath, lowResImagePath, 640, 480)
	if err != nil {
		log.Printf("failed to convert %v to %v because %v; stderr=%v", highResImagePath, lowResImagePath, err, stderr)
	}

	log.Printf("converted %v to %v", highResImagePath, lowResImagePath)
}

func (s *SegmentFolderEventHandler) Start() {
	go s.folderWatcher.Watch()

	time.Sleep(time.Second)
}

func (s *SegmentFolderEventHandler) Stop() error {
	s.folderWatcher.Stop()

	return nil
}
