# go-pttbbs

[![GoDoc](https://pkg.go.dev/badge/github.com/Ptt-official-app/go-pttbbs?status.svg)](https://pkg.go.dev/github.com/Ptt-official-app/go-pttbbs?tab=doc)
[![codecov](https://codecov.io/gh/Ptt-official-app/go-pttbbs/branch/main/graph/badge.svg)](https://codecov.io/gh/Ptt-official-app/go-pttbbs)

## README Translation

* [English](https://github.com/Ptt-official-app/go-pttbbs/blob/main/README.en.md)
* [正體中文](https://github.com/Ptt-official-app/go-pttbbs/blob/main/README.zh-TW.md)

## Overview

This project intends to be the go implementation of [ptt/pttbbs](https://github.com/ptt/pttbbs).

Collaborating with [Ptt-official-app go-pttbbsweb](https://github.com/ptt-official-app/go-pttbbsweb), go-pttbbs intends to be web-based bbs.

## Getting Started

You can start with the following steps:

* `cp docs/config/config.ini.template docs/config/02-config.dev.ini`
* update `BBSHOME` in `docs/config/02-config.dev.ini`
* `./scripts/initpasswd.sh [BBSHOME] 50`
* `./scripts/run-dev.sh`
* (in another terminal window) `./scripts/init-SYSOP.sh`

## API Document

You can start with the [swagger api](https://doc-pttbbs.devptt.dev)
and check the api document.

## Coding Convention

* [gotests](https://github.com/cweill/gotests) for test-generation
* [gofumpt](https://github.com/mvdan/gofumpt) for formatting

## Docker-Compose

You can do the following to start with docker-compose:

* copy `docs/etc/` to some etc directory (ex: `/etc/go-pttbbs`).
* copy `docs/config/01-config.docker.ini` to the etc directory as production.ini (ex: `cp docs/config/01-config.docker.ini /etc/go-pttbbs/production.ini`).
* copy `docker/go-pttbbs/docker_compose.env.template` to `docker/go-pttbbs/docker_compose.env` and modify the settings.
* `./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest`
* `docker-compose --env-file docker/go-pttbbs/docker_compose.env -f docker/go-pttbbs/docker-compose.yaml up -d`
* register SYSOP and guest (api.GUEST) at `http://localhost:3456/v1/register`
* register your account at `http://localhost:3456/register`
* login at `http://localhost:3456/v1/login`
* `telnet localhost 8888` and use the account that you registered.

## Init BBSHOME

You can do the following to init BBSHOME:

* `./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest`
* `./scripts/docker_initpasswd.sh [BBSHOME] pttofficialapps/go-pttbbs:latest [N_USER]`

## Increase Users using Docker

You can do the following to increase users using docker:

* compile new docker image (as GOPTTBBS_IMAGE) with new MAX_USER in config.go
* `./scripts/docker_tunepasswd.sh [BBSHOME] [GOPTTBBS_IMAGE]`

## Testing

```
./scripts/test.sh
```

## Coverage

```
./scripts/coverage.sh
```

## Running with run.sh

You can do the following to run with ./scripts/run.sh:

* Mac:
    Copy `docs/mac/memory.plist` to `/Library/LaunchDaemons` then reboot
    ```sh
    sudo cp docs/mac/memory.plist /Library/LaunchDaemons/memory.plist
    ```

* Check that we do have 16M shared-mem
    ```
    sysctl -a|grep shm
    ```
* Init your own BBSHOME:
    ```
    ./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest
    ```
* `cp docs/config/02-config-run.go.template ptttype/02-config-run.go`
* `cp docs/config/02-config.run.template.ini docs/config/02-config.run.ini`
* Setup BBSHOME in docs/config/02-config.run.ini
* Do the following step ONLY IF you want to reset shared-mem:
    `ipcrm -M 0x000004cc`
    `ipcrm -S 0x000007da`
* `./scripts/run.sh`

## Running in docker

You can do the following do run updated-code in docker:


* Modify the docker-compose.yaml and add the expected ports and mount directory in volumes:

    ```
    ports:
      - "127.0.0.1:3456:3456"
      - "127.0.0.1:8889:8888"
      - "127.0.0.1:48764:48763"
      - "127.0.0.1:[local-port]:[docker-port]"
    volumes:
      - ${BBSHOME}:/home/bbs
      - ${ETC}:/etc/go-pttbbs
      - [local absolute directory]:/home/[username]/go-pttbbs
    ```

* do docker-compose
* `docker container ls` and find the corresponding docker container
* `docker exec -it [container] /bin/bash`
* `cd /home/[username]/go-pttbbs`
* `./scripts/run-in-docker.sh [docker-port]`

## Config

Some config-variables are required const in ptttype,
to be defined as Cstr (IDLEN, PASSLEN, etc.)

For the normal config-variables, we use config.ini
as the configuration.

For the const config-variables in ptttype,
We use docs/config/00-config-[dev-mode].go with +build flag

### docs/config/00-config.ini
We use viper and .ini as our config-framework.
docs/config/00-config.template.ini is the config-template file.

We have 3 files For every module with the config:

1. 00-config.go: define the variables of the config.
2. config.go: define the func of setting the variables from the config-file.
3. config_util.go: helper functions. should be straightforward to follow.

### 00-config-\[dev-mode\].go

We can customized ptttype/00-config-default.go with the following steps:

1. Copy docs/config/00-config-production.go.template to ptttype/00-config-production.go and change the +build and variables accordingly.
2. `go build -tag [dev-mode]`

## Swagger.sh

The swagger setup is based on [flask-swagger](https://github.com/gangverk/flask-swagger),
which is a python-project.
You can do following for the swagger-api:

1. setup the python virtualenv.
2. cd apidoc; pip install -e . ; cd ..
3. ./scripts/swagger.sh
4. browse to [http://localhost:8080](http://localhost:8080).
