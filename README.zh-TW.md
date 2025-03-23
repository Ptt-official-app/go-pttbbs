# go-pttbbs

[![GoDoc](https://pkg.go.dev/badge/github.com/Ptt-official-app/go-pttbbs?status.svg)](https://pkg.go.dev/github.com/Ptt-official-app/go-pttbbs?tab=doc)
[![codecov](https://codecov.io/gh/Ptt-official-app/go-pttbbs/branch/main/graph/badge.svg)](https://codecov.io/gh/Ptt-official-app/go-pttbbs)

## README Translation

* [English](https://github.com/Ptt-official-app/go-pttbbs/blob/main/README.en.md)
* [正體中文](https://github.com/Ptt-official-app/go-pttbbs/blob/main/README.zh-TW.md)

## 概觀

這是使用 go 實作的 [ptt/pttbbs](https://github.com/ptt/pttbbs)

與 [Ptt-official-app go-pttbbsweb](https://github.com/ptt-official-app/go-pttbbsweb) 一起成為 web-based BBS.

## 開始

您可以從以下步驟快速開始:

* `cd docs; tar -zxvf bbs-2025-03-22.tar.gz; cd ..`
* `docker-compose --env-file docker/go-pttbbs/docker_compose.examples.env -f docker/go-pttbbs/docker-compose.yaml up -d`
* `curl -i -H 'X-Forwarded-For: 127.0.0.1' 'http://localhost:3456/v1/board/WhoAmI/articles?max=5&desc=true'`

## API Document

您可以從 [swagger api](https://doc-pttbbs.devptt.dev) 查看 api 文件.

## Coding Convention

* [gotests](https://github.com/cweill/gotests) for test-generation
* [gofumpt](https://github.com/mvdan/gofumpt) for formatting

## Docker-Compose

您可以利用 docker-compose 開始:

* 將 `docs/etc/` 複製到另一個 `etc` 目錄 (例: `/etc/go-pttbbs`).
* 將 `docs/config/01-config.docker.ini` 複製到 `etc` 目錄, 並命名為 `production.ini` (例: `cp docs/config/01-config.docker.ini /etc/go-pttbbs/production.ini`).
* 將 `docker/go-pttbbs/docker_compose.tmpl.env` 複製到 `docker/go-pttbbs/docker_compose.env` 並且更改相關設定.
* `docker-compose --env-file docker/go-pttbbs/docker_compose.env -f docker/go-pttbbs/docker-compose.yaml up -d`

## 測試

```
./scripts/test.sh
```

## Coverage

```
./scripts/coverage.sh
```

## 執行 run-dev.sh

您可以使用以下方式快速執行 `run-dev.sh`:

* `cd docs; tar -zxvf bbs-2025-03-22.tar.gz; cd ..`
* `./scripts/run-dev.sh`
* `curl -i -H 'X-Forwarded-For: 127.0.0.1' 'http://localhost:3456/v1/board/WhoAmI/articles?max=5&desc=true'`

## 設定

有些 config-variables 是在 ptttype 裡必要的常數 (IDLEN, PASSLEN, etc.)

對於一般的 config-variables, 我們在 config.ini 裡設定.

對於 ptttype 裡常數的 config-variables, 我們使用 docs/config/00-config-[dev-mode].go 和 +build flag.

### docs/config/00-config.ini

我們使用 [viper](https://github.com/spf13/viper) 和 .ini 作為我們的 config-framework.
`docs/config/00-config.tmpl.ini` 是 config-template file.

對於每個 module, 我們有以下 3 個 source files 來完成 config:

1. 00-config.go: 定義 config 的 variables.
2. config.go: 定義設定 variables 的 func.
3. config_util.go: 輔助 functions.

### docs/config/00-config-\[dev-mode\].go

我們可以使用以下方式更改 ptttype/00-config-default.go:

1. 複製 `docs/config/config-dev.go` 到 `ptttype/00-config-[dev-mode].go` 並且更改 `+build` 和 variables.
2. `go build -tag [dev-mode]`

## Swagger.sh

我們根據 [flask-swagger](https://github.com/gangverk/flask-swagger) 設定 swagger.
[flask-swagger](https://github.com/gangverk/flask-swagger) 是一個 python-project.

您可以使用以下方式設定 swagger-api:

1. 設定 python virtualenv.
2. cd apidoc; pip install -e . ; cd ..
3. ./scripts/swagger.sh
4. 使用 browser 觀看 [http://localhost:8080](http://localhost:8080).
