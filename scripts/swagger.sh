#!/bin/bash

if [ "$1" == "" ]; then
    echo "usage: swagger.sh [host]"
    exit 255
fi

host=$1

apidoc/yamltojson.py apidoc/template.yaml apidoc/template.json
flaskswagger apidoc:app --host ${host} --base-path / --out-dir swagger --from-file-keyword=swagger_from_file --template ./apidoc/template.json

docker container stop swagger-go-pttbbs
docker container rm swagger-go-pttbbs
docker run -itd --restart always --name swagger-go-pttbbs -p 8080:8080 -e SWAGGER_JSON=/foo/swagger.json -v ${PWD}/swagger:/foo swaggerapi/swagger-ui
