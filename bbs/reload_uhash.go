package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/Ptt-official-app/go-pttbbs/ptt"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func ReloadUHash(userID *ptttype.UserID_t) (err error) {
	userLevel, err := ptt.GetUserLevel(userID)
	if userLevel.HasUserPerm(ptttype.PERM_SYSOP) || err != nil {
		return api.ErrInvalidUser
	}
	return cache.LoadUHash()
}
