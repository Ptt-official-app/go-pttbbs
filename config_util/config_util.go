package config_util

import (
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/types/ansi"
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

func SetBoolConfig(configPrefix string, idx string, orig bool) bool {
	idx = configPrefix + "." + strings.ToLower(idx)
	if !viper.IsSet(idx) {
		return orig
	}

	return viper.GetBool(idx)
}

func SetColorConfig(configPrefix string, idx string, orig string) string {
	idx = configPrefix + "." + strings.ToLower(idx)
	if !viper.IsSet(idx) {
		return orig
	}
	return ansi.ANSIColor(viper.GetString(idx))
}

func SetIntConfig(configPrefix string, idx string, orig int) int {
	idx = configPrefix + "." + strings.ToLower(idx)
	if !viper.IsSet(idx) {
		return orig
	}
	return viper.GetInt(idx)
}

func SetDoubleConfig(configPrefix string, idx string, orig float64) float64 {
	idx = configPrefix + "." + strings.ToLower(idx)
	if !viper.IsSet(idx) {
		return orig
	}

	return viper.GetFloat64(idx)
}
