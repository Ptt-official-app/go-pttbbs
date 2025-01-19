package path

import "github.com/Ptt-official-app/go-pttbbs/configutil"

const configPrefix = "go-pttbbs:cmbbs:path"

func InitConfig() error {
	config()
	return nil
}

func setBoolConfig(idx string, orig bool) bool {
	return configutil.SetBoolConfig(configPrefix, idx, orig)
}
