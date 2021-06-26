#!/bin/bash

if [ "$#" != "3" ]; then
    echo "usage: docker_initpasswd.sh [BBSHOME] [GOPTTBBS_IMAGE] [N_USER]"
    exit 255
fi

BBSHOME=$1
GOPTTBBS_IMAGE=$2
N_USER=$3
echo "BBSHOME: ${BBSHOME} GOPTTBBS_IMAGE: ${GOPTTBBS_IMAGE} N_USER: ${N_USER}"

project=go-pttbbs-initpasswd

bytes=`expr ${N_USER} \* 512`
docker container stop ${project}
docker container rm ${project}
docker run --rm --name ${project} -v ${BBSHOME}:/home/bbs ${GOPTTBBS_IMAGE} dd if=/dev/zero of=/home/bbs/.PASSWDS bs=${bytes} count=1
