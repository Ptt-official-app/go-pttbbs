package ptt

import (
	"bufio"
	"bytes"
	"os"
	"strconv"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func bakumanMakeTagFilename(userID *ptttype.UserID_t, obj types.Cstr, objectType string, isCreateFile bool) (filename string, err error) {
	filename, err = path.SetHomeFile(userID, ptttype.FN_BANNED)
	if err != nil {
		return "", err
	}

	if isCreateFile && !types.DashD(filename) {
		err = types.Mkdir(filename)
		if err != nil {
			return "", err
		}
	}

	filename += string(os.PathSeparator) + objectType + "_" + types.CstrToString(obj)

	return filename, nil

}

func bakumanGetInfo(filename string) (expireTS types.Time4, reason string, err error) {
	logrus.Infof("bakumanGetInfo: start: filename: %v", filename)
	file, err := os.Open(filename)
	if err != nil {
		return 0, "", err
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	//expireTS
	line, err := types.ReadLine(buf)
	logrus.Infof("bakumanGetInfo: after read 1st line: line: %v e: %v", line, err)
	if err != nil {
		return 0, "", err
	}
	line = bytes.TrimSpace(line)
	theTS, _ := strconv.Atoi(string(line))
	expireTS = types.Time4(theTS)

	logrus.Infof("bakumanGetInfo: expireTS: %v line: %v", expireTS, string(line))

	//reason
	line, err = types.ReadLine(buf)
	if err == nil {
		line = bytes.TrimSpace(line)
		reason = string(line)
	}

	logrus.Infof("bakumanGetInfo: expireTS: %v reason: %v", expireTS, reason)

	return expireTS, reason, nil
}

//isBannedByBoard
//
//XXX TODO: implement details.
func isBannedByBoard(user *ptttype.UserecRaw, board *ptttype.BoardHeaderRaw) (expireTS types.Time4, reason string) {
	return isBannedBy(&user.UserID, board.Brdname[:], BAKUMAN_OBJECT_TYPE_BOARD)
}

func isBannedBy(userID *ptttype.UserID_t, obj types.Cstr, objectType string) (expireTS types.Time4, reason string) {

	tagFn, err := bakumanMakeTagFilename(userID, obj, objectType, false)
	if err != nil {
		return 0, ""
	}

	now := types.NowTS()
	expireTS, reason, err = bakumanGetInfo(tagFn)
	if err == nil && now > expireTS {
		_ = os.Remove(tagFn)
		return 0, ""
	}

	return expireTS, reason
}
