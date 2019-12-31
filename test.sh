#!/usr/bin/env bash

set -e

go test ./internal/api
go test ./internal/common
go test ./internal/event_renderer
go test ./internal/event_store
go test ./internal/file_converter
go test ./internal/file_diff_file_watcher
go test ./internal/file_watcher
go test ./internal/file_write_folder_watcher
go test ./internal/motion_config
go test ./internal/motion_log
go test ./internal/segment_recorder
go test ./internal/thumbnail_creator

go test ./pkg/event_api
go test ./pkg/motion_config_segment_recorder
go test ./pkg/motion_log_event_handler
go test ./pkg/segment_folder_event_handler
