package types

import (
	"time"

	configutil "github.com/Ptt-official-app/go-pttbbs/configutil"
)

const configPrefix = "go-pttbbs:types"

func InitConfig() error {
	config()
	return postConfig()
}

func setStringConfig(idx string, orig string) string {
	return configutil.SetStringConfig(configPrefix, idx, orig)
}

func setIntConfig(idx string, orig int) int {
	return configutil.SetIntConfig(configPrefix, idx, orig)
}

func setBoolConfig(idx string, orig bool) bool {
	return configutil.SetBoolConfig(configPrefix, idx, orig)
}

func postConfig() (err error) {
	_, err = setTimeLocation(TIME_LOCATION)
	if err != nil {
		return err
	}

	err = initBig5()
	if err != nil {
		return err
	}

	return nil
}

//setTimeLocation
//
//
func setTimeLocation(timeLocation string) (origTimeLocation string, err error) {
	origTimeLocation = TIME_LOCATION
	TIME_LOCATION = timeLocation

	TIMEZONE, err = time.LoadLocation(TIME_LOCATION)

	return origTimeLocation, err
}
