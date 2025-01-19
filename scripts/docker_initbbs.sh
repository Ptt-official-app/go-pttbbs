#!/bin/bash
# Initializing the BBSHOME
#
# 1. attach BBSHOME to the docker container
# 2. use /opt/bbs/bin/initbbs -DoIt to initialize .PASSWDS, home/, .BRD, boards/, and man/
#    https://github.com/ptt/pttbbs/blob/master/util/initbbs.c
# 3. copy bin/, etc/, and wsproxy/

if [ "$#" != "2" ]; then
    echo "usage: docker_initbbs.sh [BBSHOME] [GOPTTBBS_IMAGE]"
    exit 255
fi

BBSHOME=$1
GOPTTBBS_IMAGE=$2
echo "BBSHOME: ${BBSHOME} GOPTTBBS_IMAGE: ${GOPTTBBS_IMAGE}"

project=go-pttbbs-initbbs

docker container stop ${project}
docker container rm ${project}
docker run --rm --name ${project} -v ${BBSHOME}:/home/bbs ${GOPTTBBS_IMAGE} /opt/bbs/bin/initbbs -DoIt

docker run --rm --name ${project} -v ${BBSHOME}:/home/bbs ${GOPTTBBS_IMAGE} sh -c 'mkdir -p /home/bbs/bin && cp /opt/bbs/bin/* /home/bbs/bin'

docker run --rm --name ${project} -v ${BBSHOME}:/home/bbs ${GOPTTBBS_IMAGE} sh -c 'mkdir -p /home/bbs/etc && cp /opt/bbs/etc/* /home/bbs/etc'

docker run --rm --name ${project} -v ${BBSHOME}:/home/bbs ${GOPTTBBS_IMAGE} sh -c 'mkdir -p /home/bbs/wsproxy && cp -R /opt/bbs/wsproxy/* /home/bbs/wsproxy'
