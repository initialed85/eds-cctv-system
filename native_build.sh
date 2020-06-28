#!/usr/bin/env bash

set -e -x

rm -fr bin || true
mkdir -p bin || true

unset GOPATH
unset GOROOT

go build -o bin/motion_config_segment_recorder cmd/motion_config_segment_recorder/main.go
go build -o bin/motion_log_event_handler cmd/motion_log_event_handler/main.go
go build -o bin/segment_folder_event_handler cmd/segment_folder_event_handler/main.go
go build -o bin/event_store_updater_page_renderer cmd/event_store_updater_page_renderer/main.go
go build -o bin/static_file_server cmd/static_file_server/main.go
go build -o bin/event_store_deduplicator cmd/event_store_deduplicator/main.go
