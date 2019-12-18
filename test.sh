#!/usr/bin/env bash

set -e

go test ./internal/file_watcher
go test ./internal/event_streamer
go test ./internal/thumbnail_creator
go test ./internal/segment_recorder
go test ./internal/common
go test ./internal/motion_log
go test ./internal/motion_config
go test ./internal/motion_config
go test ./internal/folder_watcher
go test ./internal/file_converter

go test ./pkg/event_store
go test ./pkg/segment_folder_event_handler
go test ./pkg/motion_log_event_handler
go test ./pkg/motion_config_segment_recorder
