#!/bin/bash

ini_filename=02-config.run.ini
package=github.com/Ptt-official-app/go-pttbbs/types
commit=`git rev-parse --short HEAD`
version=`git describe --tags`

echo "commit: ${commit} version: ${version}"

go build -ldflags "-X ${package}.GIT_VERSION=${commit} -X ${package}.VERSION=${version}" -tags dev && ./go-pttbbs -ini ${ini_filename}
