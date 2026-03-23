#!/bin/bash
# build docker.

# branch
branch=`git rev-parse --abbrev-ref HEAD`
if [ "${branch}" == "HEAD" ]; then branch=`git describe --tags`; fi

# project
project=`basename \`pwd\``

docker buildx build --platform linux/amd64,linux/arm64 -t pttofficialapps/${project}:${branch} -f docker/Dockerfile --push .
