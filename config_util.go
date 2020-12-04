package main

import "github.com/Ptt-official-app/go-pttbbs/config_util"

const configPrefix = "go-pttbbs"

func InitConfig() error {
	config()
	return nil
}

func setStringConfig(idx string, orig string) string {
	return config_util.SetStringConfig(configPrefix, idx, orig)
}
