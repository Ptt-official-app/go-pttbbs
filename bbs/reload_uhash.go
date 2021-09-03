package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func ReloadUHash(uuserID UUserID) (err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return err
	}

	_, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return err
	}

	if !userecRaw.UserLevel.HasUserPerm(ptttype.PERM_SYSOP) {
		return ErrInvalidPermission
	}

	return cache.LoadUHash()
}
