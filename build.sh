#!/usr/bin/env bash

set -e -x

pushd "$(pwd)"
function teardown() {
  popd
}
trap teardown exit

if [[ ! -d quotanizer ]]; then
  git clone https://github.com/initialed85/quotanizer.git
fi

cd quotanizer
git reset --hard
git pull --all
cd ..

export CCTV_EVENTS_PATH=/tmp/cctv_events_path
export CCTV_SEGMENTS_PATH=/tmp/cctv_segments_path

mkdir -p ${CCTV_EVENTS_PATH}
mkdir -p ${CCTV_SEGMENTS_PATH}

export DOCKER_BUILDKIT=1
if ! docker-compose build; then
  echo "error: build failed"
  exit 1
fi
