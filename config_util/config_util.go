package config_util

import (
	"strings"

	"github.com/spf13/viper"
)

func SetStringConfig(configPrefix string, idx string, orig string) string {
	idx = configPrefix + "." + strings.ToLower(idx)
	if !viper.IsSet(idx) {
		return orig
	}

	return viper.GetString(idx)
}

func SetBytesConfig(configPrefix string, idx string, orig []byte) []byte {
	idx = configPrefix + "." + strings.ToLower(idx)
	if !viper.IsSet(idx) {
		return orig
	}

	return []byte(viper.GetString(idx))
}
