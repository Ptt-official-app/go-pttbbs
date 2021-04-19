#!/bin/bash

docker-compose --env-file docker_compose.env -f docker-compose.yaml  run go-pttbbs -- cp '/opt/bbs/bin/*' /home/bbs/bin
docker-compose --env-file docker_compose.env -f docker-compose.yaml up -d
