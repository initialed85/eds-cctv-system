#!/usr/bin/env bash

set -e

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

export DOCKER_BUILDKIT=1
if ! docker-compose build; then
  echo "error: build failed"
  exit 1
fi
