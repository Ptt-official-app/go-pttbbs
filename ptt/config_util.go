package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
)

func InitConfig() error {
	return postInitConfig()
}

func postInitConfig() error {
	err := cmbbs.PasswdInit()
	if err != nil {
		return err
	}

	return nil
}
