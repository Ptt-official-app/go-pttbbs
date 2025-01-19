#!/bin/bash

branch=`git rev-parse --abbrev-ref HEAD`
if [ "${branch}" == "HEAD" ]; then branch=`git describe --tags`; fi
project=`basename \`pwd\``

BBSHOME=${1:-${BBSHOME}}

docker container stop ${project}
docker container rm ${project}

echo "BBSHOME: ${BBSHOME}"
if [ "${BBSHOME}" == "" ]; then
    docker run -itd --name ${project} -p 3456:3456 -p 8888:8888 -p 48763:48763 ${project}:${branch}
else
    docker run -itd --name ${project} -p 3456:3456 -p 8888:8888 -p 48763:48763 -v ${BBSHOME}:/home/bbs ${project}:${branch}
fi
