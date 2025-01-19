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

ini_filename_prefix=`echo "${current_dir}"|awk -F '\/' '{print $3}'`
ini_filename="${ini_filename_prefix}.ini"
echo "ini_filename: ${ini_filename}"

sed "s/^HTTP_HOST =.*/HTTP_HOST = 0.0.0.0:${port}/g" /etc/go-pttbbs/production.ini > /etc/go-pttbbs/${ini_filename}

cp /srv/go-pttbbs/ptttype/00-config-production.go ptttype/00-config-production.go

go build -tags production && sudo -iu bbs ${current_dir}/go-pttbbs -ini ${ini_filename}
