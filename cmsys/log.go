package cmsys

import (
	"os"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func LogFilef(filename string, logFlag LogFlag, msg string) (err error) {
	return LogFile(filename, logFlag, msg)
}

func LogFile(filename string, logFlag LogFlag, msg string) (err error) {
	flag := os.O_APPEND | os.O_WRONLY
	if logFlag&LOG_CREAT != 0 {
		flag |= os.O_CREATE
	}

	file, err := os.OpenFile(filename, flag, ptttype.DEFAULT_FILE_CREATE_PERM)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(msg)
	if err != nil {
		return err
	}

	return nil
}
