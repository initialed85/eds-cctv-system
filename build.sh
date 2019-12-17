#!/usr/bin/env bash

source env.sh

set -e

rm -fr bin || true
mkdir -p bin || true

go build -x -v -o motion_config_segment_recorder cmd/motion_config_segment_recorder/main.go
go build -x -v -o motion_log_event_handler cmd/motion_log_event_handler/main.go
go build -x -v -o segment_folder_event_handler cmd/segment_folder_event_handler/main.go
