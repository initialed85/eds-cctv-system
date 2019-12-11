#!/usr/bin/env bash

source env.sh

docker run --rm -it --name "${CONTAINER_NAME}" "${IMAGE_NAME}"
