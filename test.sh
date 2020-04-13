#!/usr/bin/env bash

source env.sh

set -e

IMAGE_NAME="${IMAGE_NAME}-test"
CONTAINER_NAME="${CONTAINER_NAME}-test"

export DOCKER_BUILDKIT=1
docker build -t "${IMAGE_NAME}" .
docker run --rm -t --name "${CONTAINER_NAME}" "${IMAGE_NAME}"

exit 1
