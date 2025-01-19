#!/bin/bash
# build with tag=production and run go-pttbbs.

if [ "$#" -lt "1" ]
then
    echo "usage: run-production.sh [ini-filename]"
    exit 255
fi

tags=production
ini_filename=$1
types_package=github.com/Ptt-official-app/go-pttbbs/types
commit=`git rev-parse --short HEAD`
version=`git describe --tags`

echo "to build: tags: ${tags} ini: %{ini_filename} commit: ${commit} version: ${version}"

go build -ldflags "-X ${types_package}.GIT_VERSION=${commit} -X ${types_package}.VERSION=${version}" -tags ${tags}
echo "to run go-pttbbs"
go-pttbbs -ini ${ini_filename}
