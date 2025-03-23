# go-pttbbs

[![GoDoc](https://pkg.go.dev/badge/github.com/Ptt-official-app/go-pttbbs?status.svg)](https://pkg.go.dev/github.com/Ptt-official-app/go-pttbbs?tab=doc)
[![codecov](https://codecov.io/gh/Ptt-official-app/go-pttbbs/branch/main/graph/badge.svg)](https://codecov.io/gh/Ptt-official-app/go-pttbbs)

## README Translation

- [English](./README.en.md)
- [正體中文](./README.zh-TW.md)

## Overview

This project intends to be the go implementation of [ptt/pttbbs](https://github.com/ptt/pttbbs).

Collaborating with [Ptt-official-app go-pttbbsweb](https://github.com/ptt-official-app/go-pttbbsweb), go-pttbbs intends to be web-based bbs.

## Getting Started

You can start with the following steps:

- `cd docs; tar -zxvf bbs-2025-03-22.tar.gz; cd ..`
- `docker-compose --env-file docker/go-pttbbs/docker_compose.examples.env -f docker/go-pttbbs/docker-compose.yaml up -d`
- `curl -i -H 'X-Forwarded-For: 127.0.0.1' 'http://localhost:3456/v1/board/WhoAmI/articles?max=5&desc=true'`

## API Document

You can start with the [swagger api](https://doc-pttbbs.devptt.dev)
and check the api document.

## Coding Convention

- [gotests](https://github.com/cweill/gotests) for test-generation
- [gofumpt](https://github.com/mvdan/gofumpt) for formatting

## Docker-Compose

You can do the following to start with docker-compose:

- copy `docs/etc/` to some `etc` directory (ex: `/etc/go-pttbbs`).
- copy `docs/config/01-config.docker.ini` to the `etc` directory as `production.ini` (ex: `cp docs/config/01-config.docker.ini /etc/go-pttbbs/production.ini`).
- copy `docker/go-pttbbs/docker_compose.tmpl.env` to `docker/go-pttbbs/docker_compose.env` and modify the settings.
- `docker-compose --env-file docker/go-pttbbs/docker_compose.env -f docker/go-pttbbs/docker-compose.yaml up -d`

## Testing

```sh
./scripts/test.sh
```

## Coverage

```sh
./scripts/coverage.sh
```

## run-dev.sh

You can do the following to execute `./scripts/run.sh`:

- `cd docs; tar -zxvf bbs-2025-03-22.tar.gz; cd ..`
- `./scripts/run-dev.sh`
- `curl -i -H 'X-Forwarded-For: 127.0.0.1' 'http://localhost:3456/v1/board/WhoAmI/articles?max=5&desc=true'`

## Config

Some config-variables are required const in ptttype,
to be defined as Cstr (IDLEN, PASSLEN, etc.)

For the normal config-variables, we use config.ini
as the configuration.

For the const config-variables in ptttype,
We use docs/config/00-config-[dev-mode].go with +build flag

### docs/config/00-config.ini

We use [viper](https://github.com/spf13/viper) and .ini as our config-framework.
`docs/config/config.tmpl.ini` is the config-template file.

We have 3 files For every module with the config:

1. `00-config.go`: define the variables of the config.
2. `config.go`: define the func of setting the variables from the config-file.
3. `config_util.go`: helper functions. should be straightforward to follow.

### 00-config-\[dev-mode\].go

We can customized `ptttype/00-config-default.go` with the following steps:

1. copy `docs/config/config-dev.go` to `ptttype/00-config-[dev-mode].go` and change the `+build` and variables accordingly.
2. `go build -tag [dev-mode]`

## Swagger.sh

The swagger setup is based on [flask-swagger](https://github.com/gangverk/flask-swagger),
which is a python-project.
You can do following for the swagger-api:

1. setup the python virtualenv.
2. cd apidoc; pip install -e . ; cd ..
3. ./scripts/swagger.sh
4. browse to [http://localhost:8080](http://localhost:8080).
