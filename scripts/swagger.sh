#!/bin/bash
# generating swagger api document.

if [ "$1" == "" ]; then
    echo "usage: swagger.sh [host]"
    exit 255
fi

host=$1

apidoc/yamltojson.py apidoc/template.yaml apidoc/template.json
flaskswagger apidoc --host ${host} --base-path / --out-dir swagger --from-file-keyword=swagger_from_file --template ./apidoc/template.json
docker run --rm -v ${PWD}/swagger:/data swaggerapi/swagger-codegen-cli-v3 generate -l openapi-yaml -i file:///data/swagger.json -o /data/v3
apidoc/addsecurity.py swagger/v3/openapi.yaml
apidoc/parameters.py swagger/swagger.json swagger/v3/openapi.security.json

docker container stop swagger-go-pttbbs
docker container rm swagger-go-pttbbs
docker run -itd --restart always --name swagger-go-pttbbs -p 127.0.0.1:8080:8080 -e SWAGGER_JSON=/foo/v3/openapi.security.params.json -v ${PWD}/swagger:/foo swaggerapi/swagger-ui
