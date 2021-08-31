package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func ReloadUHash(userID *ptttype.UserID_t) (err error) {
	userLevel, err := ptt.GetUserLevel(userID)
	if err != nil {
		return err
	}
	if userLevel.HasUserPerm(ptttype.PERM_SYSOP) {
		return cache.LoadUHash()
	}
	return ErrInvalidPermission
}
