#!/usr/bin/env bash

set -e -x

unset GOPATH
unset GOROOT

go test -v ./internal/api
go test -v ./internal/common
go test -v ./internal/duration_finder
go test -v ./internal/page_renderer
go test -v ./internal/event_store
go test -v ./internal/event_store_updater
go test -v ./internal/file_converter
go test -v ./internal/file_diff_file_watcher
go test -v ./internal/file_watcher
go test -v ./internal/file_write_folder_watcher
go test -v ./internal/motion_config
go test -v ./internal/motion_log
go test -v ./internal/segment_recorder
go test -v ./internal/thumbnail_creator

go test -v ./pkg/event_api
go test -v ./pkg/event_store_updater_page_renderer
go test -v ./pkg/motion_config_segment_recorder
go test -v ./pkg/motion_log_event_handler
go test -v ./pkg/segment_folder_event_handler
go test -v ./pkg/event_store_deduplicator

# no venv in the test container
source venv/bin/activate || true
python3 -m pytest -v -s utils
