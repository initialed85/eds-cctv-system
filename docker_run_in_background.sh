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
  -p 81:80 \
  -p 9080:8080 \
  -p 9081:8081 \
  -p 9082:8082 \
  -p 9083:8083 \
  -p 9084:8084 \
  -v "$(pwd)"/motion-configs:/etc/motion \
  -v /media/storage/Cameras/events:/srv/target_dir/events \
  -v /media/storage/Cameras/segments:/srv/target_dir/segments \
  "${IMAGE_NAME}"
