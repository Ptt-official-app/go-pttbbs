#!/bin/bash

ini_filename=00-config.template.ini

echo "to build"
go build

echo "to run go-pttbbs"
./go-pttbbs -ini ${ini_filename}
