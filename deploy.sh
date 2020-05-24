#!/usr/bin/env bash

set -e -x

pushd "$(pwd)"
function teardown() {
  docker-compose down || true
}
trap teardown ERR

if [[ "${CCTV_MOTION_CONFIGS}" == "" ]]; then
  echo "error: CCTV_MOTION_CONFIGS environment variable not set"

  exit 1
fi

if [[ "${CCTV_EVENTS_PATH}" == "" ]]; then
  echo "error: CCTV_EVENTS_PATH environment variable not set"

  exit 1
fi

if [[ "${CCTV_SEGMENTS_PATH}" == "" ]]; then
  echo "error: CCTV_SEGMENTS_PATH environment variable not set"

  exit 1
fi

if [[ "${CCTV_EVENTS_QUOTA}" == "" ]]; then
  echo "error: CCTV_EVENTS_QUOTA environment variable not set"

  exit 1
fi

if [[ "${CCTV_SEGMENTS_QUOTA}" == "" ]]; then
  echo "error: CCTV_SEGMENTS_QUOTA environment variable not set"

  exit 1
fi

export CCTV_MOTION_CONFIGS="${CCTV_MOTION_CONFIGS}"
export CCTV_EVENTS_PATH="${CCTV_EVENTS_PATH}"
export CCTV_SEGMENTS_PATH="${CCTV_SEGMENTS_PATH}"
export CCTV_EVENTS_QUOTA="${CCTV_EVENTS_QUOTA}"
export CCTV_SEGMENTS_QUOTA="${CCTV_SEGMENTS_QUOTA}"

export DOCKER_BUILDKIT=1
docker-compose down || true
docker-compose up -d
