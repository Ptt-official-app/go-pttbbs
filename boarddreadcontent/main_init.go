package main

import (
	"flag"
	"strings"

	"github.com/spf13/viper"

	"github.com/Ptt-official-app/go-pttbbs/boardd"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
	jww "github.com/spf13/jwalterweatherman"
)

//Params
//      filename: ini filename
//
//Return
//      error: err
func initAllConfig(filename string) error {
	filenameList := strings.Split(filename, ".")
	if len(filenameList) == 1 {
		return ErrInvalidIni
	}

	filenamePrefix := strings.Join(filenameList[:len(filenameList)-1], ".")
	filenamePostfix := filenameList[len(filenameList)-1]
	viper.SetConfigName(filenamePrefix)
	viper.SetConfigType(filenamePostfix)
	viper.AddConfigPath("/etc/go-pttbbs")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	logrus.Debugf("viper keys: %v", viper.AllKeys())

	err = types.InitConfig()
	if err != nil {
		return err
	}

	err = boardd.InitConfig()
	if err != nil {
		return err
	}

	return nil
}

func initMain() (filename string, err error) {
	jww.SetLogThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelDebug)
	logrus.SetLevel(logrus.InfoLevel)

	iniFilename := ""
	flag.StringVar(&iniFilename, "ini", "config.ini", "ini filename")

	flag.StringVar(&filename, "file", "", "filename")

	flag.Parse()

	err = initAllConfig(iniFilename)
	if err != nil {
		return "", err
	}

	err = boardd.Init(false)
	if err != nil {
		return "", err
	}

	return filename, nil
}
