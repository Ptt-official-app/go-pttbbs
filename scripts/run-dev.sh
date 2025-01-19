#!/bin/bash
# build with tag=dev and run go-pttbbs with 02-config.run.ini.

tags=dev
ini_filename=docs/config/02-config.dev.ini
types_package=github.com/Ptt-official-app/go-pttbbs/types
commit=`git rev-parse --short HEAD`
version=`git describe --tags`

echo "commit: ${commit} version: ${version}"

cd go-pttbbs
go build -ldflags "-X ${types_package}.GIT_VERSION=${commit} -X ${types_package}.VERSION=${version}" -tags ${tags}
cd ..

./go-pttbbs/go-pttbbs -ini ${ini_filename}
