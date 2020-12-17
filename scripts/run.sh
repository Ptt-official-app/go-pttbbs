#!/bin/bash

ini_filename=02-config.run.ini

echo "to build"
go build -tags dev

echo "to run go-pttbbs"
./go-pttbbs -ini ${ini_filename}
