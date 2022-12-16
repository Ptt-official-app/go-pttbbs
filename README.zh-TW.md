# go-pttbbs

[![GoDoc](https://pkg.go.dev/badge/github.com/Ptt-official-app/go-pttbbs?status.svg)](https://pkg.go.dev/github.com/Ptt-official-app/go-pttbbs?tab=doc)
[![codecov](https://codecov.io/gh/Ptt-official-app/go-pttbbs/branch/main/graph/badge.svg)](https://codecov.io/gh/Ptt-official-app/go-pttbbs)

## README Translation

* [English](https://github.com/Ptt-official-app/go-pttbbs/blob/main/README.en.md)
* [正體中文](https://github.com/Ptt-official-app/go-pttbbs/blob/main/README.zh-TW.md)

## 概觀

這是使用 go 實作的 [ptt/pttbbs](https://github.com/ptt/pttbbs)

與 [Ptt-official-app middleware](https://github.com/ptt-official-app/go-openbbsmiddleware) 一起成為 web-based BBS.

## 開始

您可以從 [swagger api](https://api.devptt.dev:8080) 查看 api 文件.

## Coding Convention

* [gotests](https://github.com/cweill/gotests) for test-generation
* [gofumpt](https://github.com/mvdan/gofumpt) for formatting

## Docker-Compose

您可以利用 docker-compose 開始:

* 將 `docs/etc/` 複製到另一個 etc 目錄 (例: `/etc/go-pttbbs`).
* 將 `01-config.docker.ini` 複製到 etc 目錄為 production.ini (例: `cp 01-config.docker.ini /etc/go-pttbbs/production.ini`).
* 將 `docker_compose.env.template` 複製到 `docker_compose.env` 並且更改相關設定.
* `./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest`
* `docker-compose --env-file docker_compose.env -f docker-compose.yaml up -d`
* 使用 `http://localhost:3456/v1/register` 註冊 SYSOP 和 guest (api.GUEST).
* 使用 `http://localhost:3456/register` 註冊您的帳號.
* 使用 `http://localhost:3456/v1/login` 登入.
* `telnet localhost 8888` 並且使用您的帳號登入.

## 起始 BBSHOME

您可以使用以下方式起始 BBSHOME:

* `./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest`
* `./scripts/docker_initpasswd.sh [BBSHOME] pttofficialapps/go-pttbbs:latest [N_USER]`

## 使用 Docker 增加使用者.

您可以使用以下方式增加使用者:

* 在 config.go 設定新的 `MAX_USER`, 並編譯成新的 docker image (as GOPTTBBS_IMAGE).
* `./scripts/docker_tunepasswd.sh [BBSHOME] [GOPTTBBS_IMAGE]`

## 測試

```
./scripts/test.sh
```

## Coverage

```
./scripts/coverage.sh
```

## 執行 run.sh

您可以使用以下方式執行 ./scripts/run.sh

* Mac:
    將 `memory.plist` 複製到 `/Library/LaunchDaemons` 然後重新啟動 Mac.
    ```sh
    sudo cp memory.plist /Library/LaunchDaemons/memory.plist
    ```

* 確認我們真的有 16M shared-mem
    ```
    sysctl -a|grep shm
    ```
* 起始您的 BBSHOME:
    ```
    ./scripts/docker_initbbs.sh [BBSHOME] pttofficialapps/go-pttbbs:latest
    ```
* `cp 02-config-run.go.template ptttype/02-config-run.go`
* `cp 02-config.run.template.ini 02-config.run.ini`
* 在 02-config.run.ini 裡設定 BBSHOME.
* "只有如果" 您想要重新設定 shared-mem, 您可以使用以下方式:
    `ipcrm -M 0x000004cc`
    `ipcrm -S 0x000007da`
* `./scripts/run.sh`

## 在 Docker 裡執行

您可以使用以下方式在 Docker 裡執行:

* 更改 docker-compose.yaml 裡的設定, 並且加上相對應的 ports 和目錄:

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

* 執行 [docker-compose](#Docker-Compose)
* `docker container ls` 並找到相對應的 docker-container
* `docker exec -it [container] /bin/bash`
* `cd /home/[username]/go-pttbbs`
* `./scripts/run-in-docker.sh [docker-port]`

## 設定

有些 config-variables 是在 ptttype 裡必要的常數 (IDLEN, PASSLEN, etc.)

對於一般的 config-variables, 我們在 config.ini 裡設定.

對於 ptttype 裡常數的 config-variables, 我們使用 00-config-[dev-mode].go 和 +build flag.

### 00-config.ini

我們使用 viper 和 .ini 作為我們的 config-framework.
00-config.template.ini 是 config-template file.

對於每個 module, 我們有以下 3 個 source files 來完成 config:

1. 00-config.go: 定義 config 的 variables.
2. config.go: 定義設定 variables 的 func.
3. config_util.go: 輔助 functions.

### 00-config-\[dev-mode\].go

我們可以使用以下方式更改 ptttype/00-config-default.go:

1. 複製 00-config-production.go.template 到 ptttype/00-config-[dev-mode].go 並且更改 +build 和 variables.
2. `go build -tag [dev-mode]`

## Swagger.sh

我們根據 [flask-swagger](https://github.com/gangverk/flask-swagger) 設定 swagger.
[flask-swagger](https://github.com/gangverk/flask-swagger) 是一個 python-project.

您可以使用以下方式設定 swagger-api:

1. 設定 python virtualenv.
2. cd apidoc; pip install . && pip uninstall apidoc -y && python setup.py develop; cd ..
3. ./scripts/swagger.sh
4. 使用 browser 觀看 [http://localhost:8080](http://localhost:8080).
