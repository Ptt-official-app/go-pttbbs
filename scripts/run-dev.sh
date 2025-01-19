#!/bin/bash
# build with tag=dev and run go-pttbbs with 02-config.run.ini.

ini_filename=docs/config/02-config.dev.ini
package=github.com/Ptt-official-app/go-pttbbs/types
commit=`git rev-parse --short HEAD`
version=`git describe --tags`

echo "commit: ${commit} version: ${version}"

cd go-pttbbs
go build -ldflags "-X ${package}.GIT_VERSION=${commit} -X ${package}.VERSION=${version}" -tags dev 
cd ..

./go-pttbbs/go-pttbbs -ini ${ini_filename}
