package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func SetUserPerm(userID UUserID, setUserID UUserID, perm ptttype.PERM) (newPerm ptttype.PERM, err error) {
	userIDRaw, err := userID.ToRaw()
	if err != nil {
		return ptttype.PERM_INVALID, ErrInvalidParams
	}

	userecRaw, err := ptt.GetUser(userIDRaw)
	if err != nil {
		return ptttype.PERM_INVALID, err
	}

	setUserIDRaw, err := setUserID.ToRaw()
	if err != nil {
		return ptttype.PERM_INVALID, ErrInvalidParams
	}

	setUID, setUserecRaw, err := ptt.InitCurrentUser(setUserIDRaw)
	if err != nil {
		return ptttype.PERM_INVALID, err
	}

	return ptt.SetUserPerm(userecRaw, setUID, setUserecRaw, perm)
}
