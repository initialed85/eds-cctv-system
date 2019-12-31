#!/usr/bin/env bash

source env.sh

set -e

rm -fr bin || true
mkdir -p bin || true

go build -v -o bin/motion_config_segment_recorder cmd/motion_config_segment_recorder/main.go
go build -v -o bin/motion_log_event_handler cmd/motion_log_event_handler/main.go
go build -v -o bin/segment_folder_event_handler cmd/segment_folder_event_handler/main.go
go build -v -o bin/event_store_updater_event_renderer cmd/event_store_updater_event_renderer/main.go
