#!/bin/bash

ini_filename=02-config.run.ini
package=github.com/Ptt-official-app/go-pttbbs/types
commit=$(shell git rev-parse --short HEAD)

go build -ldflags "-X ${package}.GIT_VERSION=${commit}" -tags dev && ./go-pttbbs -ini ${ini_filename}
