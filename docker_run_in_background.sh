#!/usr/bin/env bash

source env.sh

set -e

docker stop "${CONTAINER_NAME}" || true
docker rm -f "${CONTAINER_NAME}" || true
# --

docker run \
  -d \
  --restart=always \
  --name "${CONTAINER_NAME}" \
  -p 80:80 \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 8082:8082 \
  -p 8083:8083 \
  -p 8084:8084 \
  -v "$(pwd)"/motion-configs:/etc/motion \
  -v /media/storage/Cameras/events:/srv/target_dir/events \
  -v /media/storage/Cameras/segments:/srv/target_dir/segments \
  "${IMAGE_NAME}"
