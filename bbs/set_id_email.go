package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func SetIDEmail(uuserID UUserID, isSet bool) (userLevel2 ptttype.PERM2, err error) {

	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return ptttype.PERM2_INVALID, err
	}

	return ptt.ChangeUserLevel2(userIDRaw, ptttype.PERM2_ID_EMAIL, isSet)
}
