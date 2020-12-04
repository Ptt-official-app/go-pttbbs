# go-pttbbs

This project intends to be the go implementation of pttbbs/pttbbs.

Collaborating with Ptt-official-app middlewares.
go-pttbbs intends to be web-based bbs.

## Testing

```
go test ./...
```

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
