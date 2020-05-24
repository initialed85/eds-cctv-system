#!/usr/bin/env bash

set -e -x

IMAGE_NAME="eds-cctv-system-test"
CONTAINER_NAME="eds-cctv-system-test"

export DOCKER_BUILDKIT=1
docker build -t "${IMAGE_NAME}" .
docker run --rm -t --name "${CONTAINER_NAME}" "${IMAGE_NAME}"

exit 1
