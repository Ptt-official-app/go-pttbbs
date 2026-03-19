package configutil

import (
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/types/ansi"
	"github.com/go-viper/encoding/ini"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var myViper *viper.Viper

func InitViper(filename string) (err error) {
	codecRegistry := viper.NewCodecRegistry()
	_ = codecRegistry.RegisterCodec("ini", ini.Codec{})

	myViper = viper.NewWithOptions(
		viper.WithCodecRegistry(codecRegistry),
	)

	filenameList := strings.Split(filename, ".")
	if len(filenameList) == 1 {
		return ErrInvalidIni
	}

	filenamePrefix := strings.Join(filenameList[:len(filenameList)-1], ".")
	filenamePostfix := filenameList[len(filenameList)-1]

	myViper.SetConfigName(filenamePrefix)
	myViper.SetConfigType(filenamePostfix)
	myViper.AddConfigPath("/etc/go-pttbbs")
	myViper.AddConfigPath(".")
	err = myViper.ReadInConfig()
	if err != nil {
		return err
	}

	logrus.Infof("viper keys: %v", viper.AllKeys())

	return nil
}

func SetStringConfig(configPrefix string, idx string, orig string) string {
	if myViper == nil {
		myViper = viper.GetViper()
	}
	idx = configPrefix + "." + strings.ToLower(idx)
	if !myViper.IsSet(idx) {
		return orig
	}

	return myViper.GetString(idx)
}

func SetBytesConfig(configPrefix string, idx string, orig []byte) []byte {
	if myViper == nil {
		myViper = viper.GetViper()
	}
	idx = configPrefix + "." + strings.ToLower(idx)
	if !myViper.IsSet(idx) {
		return orig
	}

	return []byte(myViper.GetString(idx))
}

func SetBoolConfig(configPrefix string, idx string, orig bool) bool {
	if myViper == nil {
		myViper = viper.GetViper()
	}
	idx = configPrefix + "." + strings.ToLower(idx)
	if !myViper.IsSet(idx) {
		return orig
	}

	return myViper.GetBool(idx)
}

func SetColorConfig(configPrefix string, idx string, orig string) string {
	if myViper == nil {
		myViper = viper.GetViper()
	}
	idx = configPrefix + "." + strings.ToLower(idx)
	if !myViper.IsSet(idx) {
		return orig
	}
	return ansi.ANSIColor(myViper.GetString(idx))
}

func SetIntConfig(configPrefix string, idx string, orig int) int {
	if myViper == nil {
		myViper = viper.GetViper()
	}
	idx = configPrefix + "." + strings.ToLower(idx)
	if !myViper.IsSet(idx) {
		return orig
	}
	return myViper.GetInt(idx)
}

func SetDoubleConfig(configPrefix string, idx string, orig float64) float64 {
	if myViper == nil {
		myViper = viper.GetViper()
	}
	idx = configPrefix + "." + strings.ToLower(idx)
	if !myViper.IsSet(idx) {
		return orig
	}

	return myViper.GetFloat64(idx)
}
