#!/usr/bin/env bash

source env.sh

test -e

docker stop "${CONTAINER_NAME}" || true

docker rm -f "${CONTAINER_NAME}" || true

docker run \
  -d \
  --restart=always \
  --name "${CONTAINER_NAME}" \
  -p 80:80 \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 8082:8082 \
  -v "$(pwd)"/motion-configs:/etc/motion \
  -v /media/storage/Cameras:/srv/target_dir \
  "${IMAGE_NAME}"
