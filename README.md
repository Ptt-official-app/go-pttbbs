# go-pttbbs

This project intends to be the go implementation of pttbbs/pttbbs.

Collaborating with Ptt-official-app middlewares.
go-pttbbs intends to be web-based bbs.

## Getting Started

You can start with the [swagger api](http://173.255.216.176:8080)
and try the api.

You can copy the curl command from the link if you encounter
CORS issue.

## Docker-Compose

You can do the following to start with docker-compose:

* copy `docker_compose.env.template` to `docker_compose.env` and modify the settings.
* `./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest`
* `docker-compose --env-file docker_compose.env -f docker-compose.yaml up -d`
* register SYSOP and pttguest (api.GUEST) at `http://localhost:3456/register`
* register your account at `http://localhost:3456/register`
* login at `http://localhost:3456/v1/login`
* `telnet localhost 8888` and use the account that you registered.

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


## Relationship with https://github.com/PichuChen/go-bbs

[PiChuChen's go-bbs](https://github.com/PichuChen/go-bbs) was the first repo
in Ptt-official-app intending to be the backend. This repository is intending to have
another approach to implement the go-version of the bbs.
