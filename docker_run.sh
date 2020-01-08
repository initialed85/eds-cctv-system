#!/usr/bin/env bash

source env.sh

set -e

mkdir output || true
mkdir -p output/events || true
mkdir -p output/segments || true

docker run \
  --rm \
  -it \
  --name "${CONTAINER_NAME}" \
  -p 81:80 \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 8082:8082 \
  -p 8083:8083 \
  -p 8084:8084 \
  -v "$(pwd)"/motion-configs:/etc/motion \
  -v "$(pwd)"/output/events:/srv/target_dir/events \
  -v "$(pwd)"/output/segments:/srv/target_dir/segments \
  "${IMAGE_NAME}"
