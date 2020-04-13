#!/usr/bin/env bash

source env.sh

set -e

rm -fr bin || true
mkdir -p bin || true

unset GOPATH
unset GOROOT

go build -v -o bin/motion_config_segment_recorder cmd/motion_config_segment_recorder/main.go
go build -v -o bin/motion_log_event_handler cmd/motion_log_event_handler/main.go
go build -v -o bin/segment_folder_event_handler cmd/segment_folder_event_handler/main.go
go build -v -o bin/event_store_updater_page_renderer cmd/event_store_updater_page_renderer/main.go
go build -v -o bin/static_file_server cmd/static_file_server/main.go
go build -v -o bin/event_store_deduplicator cmd/event_store_deduplicator/main.go
