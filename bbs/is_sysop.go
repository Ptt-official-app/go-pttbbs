package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func IsSysop(uuserID UUserID, perm ptttype.PERM) (isValid bool) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return false
	}

	userLevel, err := ptt.GetUserLevel(userIDRaw)
	if err != nil {
		return false
	}

	return userLevel&perm != 0
}
