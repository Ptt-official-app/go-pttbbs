#!/bin/bash

branch=`git branch|grep '^*'|sed 's/^\* //g'|sed -E 's/^\(HEAD detached at //g'|sed -E 's/\)$//g'`
project=`basename \`pwd\``

BBSHOME=${1:-${BBSHOME}}

docker container stop ${project}
docker container rm ${project}

if [ "${BBSHOME}" == "" ]; then
    docker run -itd --name ${project} -p 3456:3456 -p 8888:8888 -p 48763:48763 ${project}:${branch}
else
    docker run -itd --name ${project} -p 3456:3456 -p 8888:8888 -p 48763:48763 -v ${BBSHOME}:/home/bbs ${project}:${branch}
fi
