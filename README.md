# go-pttbbs

This project intends to be the go implementation of pttbbs/pttbbs.

Collaborating with Ptt-official-app middlewares.
go-pttbbs intends to be web-based bbs.

## Testing

```
go test ./...
```

## Config

We use viper and .ini as our config-framework.
00-config.template.ini is the config-template file.

We have 3 files For every module with the config:

1. 00-config.go: define the variables of the config.
2. config.go: define the func of setting the variables from the config-file.
3. config_util.go: helper functions. should be straightforward to follow.
