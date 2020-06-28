#!/usr/bin/env bash

set -e -x

IMAGE_NAME="eds-cctv-system-test"
CONTAINER_NAME="eds-cctv-system-test"

DOCKER_ARGS=""
ENTRYPOINT_ARGS=""
if [[ -n "${1}" ]]; then
  DOCKER_ARGS="-i --entrypoint ${1}"
  ENTRYPOINT_ARGS="${*:2}"
fi

DOCKER_BUILDKIT=1 docker build -t "${IMAGE_NAME}" .

# shellcheck disable=SC2086
docker run --rm -t --name "${CONTAINER_NAME}" ${DOCKER_ARGS} "${IMAGE_NAME}" ${ENTRYPOINT_ARGS}

exit 1
