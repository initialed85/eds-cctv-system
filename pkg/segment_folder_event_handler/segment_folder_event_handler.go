package segment_folder_event_handler

import (
	"eds-cctv-system/internal/folder_watcher"
)

type SegmentFolderEventHandler struct {
	folderWatcher folder_watcher.FolderWatcher
}
