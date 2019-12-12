#!/usr/bin/env bash

source env.sh

if ! docker build -t "${IMAGE_NAME}" -f docker/Dockerfile .; then
  echo "error: build failed"

  exit 1
fi

if ! docker run -d --entrypoint true --name "${CONTAINER_NAME}"-build "${IMAGE_NAME}"; then
  echo "error: run failed (to get example configs)"

  docker rm -f "${CONTAINER_NAME}"

  exit 1
fi

rm -fr motion-config-examples >/dev/null 2>&1

mkdir -p motion-config-examples

docker cp "${CONTAINER_NAME}"-build:/etc/motion/examples/motion.conf motion-config-examples/motion.conf
docker cp "${CONTAINER_NAME}"-build:/etc/motion/examples/camera1.conf motion-config-examples/camera1.conf
docker cp "${CONTAINER_NAME}"-build:/etc/motion/examples/camera2.conf motion-config-examples/camera2.conf
docker cp "${CONTAINER_NAME}"-build:/etc/motion/examples/camera3.conf motion-config-examples/camera3.conf

docker rm -f "${CONTAINER_NAME}"-build
