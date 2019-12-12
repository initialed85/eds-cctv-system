#!/usr/bin/env bash

source env.sh

docker run \
  --rm \
  -d \
  --restart=always \
  --name "${CONTAINER_NAME}" \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 8082:8082 \
  -v "$(pwd)"/configs:/etc/motion \
  -v /media/storage/Cameras:/srv/target_dir \
  "${IMAGE_NAME}"
