# go-pttbbs

[![GoDoc](https://pkg.go.dev/badge/github.com/Ptt-official-app/go-pttbbs?status.svg)](https://pkg.go.dev/github.com/Ptt-official-app/go-pttbbs?tab=doc)
[![codecov](https://codecov.io/gh/Ptt-official-app/go-pttbbs/branch/main/graph/badge.svg)](https://codecov.io/gh/Ptt-official-app/go-pttbbs)

This project intends to be the go implementation of pttbbs/pttbbs.

Collaborating with Ptt-official-app middlewares.
go-pttbbs intends to be web-based bbs.

## Getting Started

You can start with the [swagger api](https://api.devptt.site:8080)
and try the api.

You can copy the curl command from the link if you encounter
CORS issue.

## Coding Convention

* [gotests](https://github.com/cweill/gotests) for test-generation
* [gofumpt](https://github.com/mvdan/gofumpt) for formatting

## Docker-Compose

You can do the following to start with docker-compose:

* copy `docs/etc/` to some etc directory (ex: `/etc/go-pttbbs`).
* copy `01-config.docker.ini` to the etc directory as production.ini (ex: `cp 01-config.docker.ini /etc/go-pttbbs/production.ini`).
* copy `docker_compose.env.template` to `docker_compose.env` and modify the settings.
* `./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest`
* `docker-compose --env-file docker_compose.env -f docker-compose.yaml up -d`
* register SYSOP and guest (api.GUEST) at `http://localhost:3456/v1/register`
* register your account at `http://localhost:3456/register`
* login at `http://localhost:3456/v1/login`
* `telnet localhost 8888` and use the account that you registered.

## Init BBS Home

You can do the following to init bbs-home:

* `./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest`
* `./scripts/docker_initpasswd.sh [BBSHOME] pttofficialapps/go-pttbbs:latest [N_USER]`

## Increase Users in Docker

You can do the following to increase users in docker:

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

* On Mac: put the following lines in /etc/sysctl.conf and reboot for 16M shared-mem:
    ```
    kern.sysv.shmmax=16777216
    kern.sysv.shmmin=1
    kern.sysv.shmmni=128
    kern.sysv.shmseg=32
    kern.sysv.shmall=4096
    ```
    
    For Mac Big Sur or after, the setting in `/etc/sysctl.conf` would be **deprecated**, therefore, we can use alternative way to set up shm.
    Copy `memory.plist` daemon to `/Library/LaunchDaemons` then reboot
    ```sh
    sudo cp memory.plist /Library/LaunchDaemons/memory.plist
    ```

* Check that we do have 16M shared-mem
    ```
    sysctl -a|grep shm
    ```
* Init your own BBSHOME:
    ```
    ./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest
    ```
* `cp 02-config-run.go.template ptttype/02-config-run.go`
* `cp 02-config.run.template.ini 02-config.run.ini`
* Setup BBSHOME in 02-config.run.ini
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
We use 00-config-[dev-mode].go with +build flag

### 00-config.ini
We use viper and .ini as our config-framework.
00-config.template.ini is the config-template file.

We have 3 files For every module with the config:

1. 00-config.go: define the variables of the config.
2. config.go: define the func of setting the variables from the config-file.
3. config_util.go: helper functions. should be straightforward to follow.

### 00-config-\[dev-mode\].go

We can customized ptttype/00-config-default.go with the following steps:

1. Copy 00-config-production.go.template to ptttype/00-config-production.go and change the +build and variables accordingly.
2. `cd go-pttbbs ; go build -tag [dev-mode]; cd ..`


## Swagger.sh

The swagger setup is based on [flask-swagger](https://github.com/gangverk/flask-swagger),
which is a python-project.
You can do following for the swagger-api:

1. setup the python virtualenv.
2. cd apidoc; pip install . && pip uninstall apidoc -y && python setup develop; cd ..
3. ./scripts/swagger.sh
4. connect to [http://localhost:8080](http://localhost:8080)
