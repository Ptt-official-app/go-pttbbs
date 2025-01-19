#!/bin/bash
# run customized go-pttbbs in docker.
#
# This is suitable only in `docker run -itd [GOPTTBBS_IMAGE] /bin/bash`.

if [ "$#" != "1" ]; then
    echo "usage: run-in-docker.sh [port]"
    exit 255
fi

port=$1
current_dir=`pwd`

tags=production
ini_filename="run-in-docker.ini"
types_package=github.com/Ptt-official-app/go-pttbbs/types
commit=`git rev-parse --short HEAD`
version=`git describe --tags`

echo "ini_filename: ${ini_filename} commit: ${commit} version: ${version}"
sed "s/^HTTP_HOST =.*/HTTP_HOST = 0.0.0.0:${port}/g" /etc/go-pttbbs/production.ini > /etc/go-pttbbs/${ini_filename}
cp /srv/go-pttbbs/ptttype/00-config-production.go ptttype/00-config-production.go

cd go-pttbbs
go build -ldflags "-X ${types_package}.GIT_VERSION=${commit} -X ${types_package}.VERSION=${version}" -tags ${tags}
cd ..

sudo -iu bbs ${current_dir}/go-pttbbs/go-pttbbs -ini ${ini_filename}
