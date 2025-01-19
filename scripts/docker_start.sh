#!/bin/bash
# Starting docker compose.
#
# 1. update /home/bbs/bin to the image version of the bin (/opt/bbs/bin).
# 2. docker compose up.

docker-compose --env-file docker/go-pttbbs/docker_compose.env -f docker/go-pttbbs/docker-compose.yaml run go-pttbbs -- cp '/opt/bbs/bin/*' /home/bbs/bin
docker-compose --env-file docker/go-pttbbs/docker_compose.env -f docker/go-pttbbs/docker-compose.yaml up -d
