package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func initCurrentUser(userID *ptttype.UserID_t) (uid ptttype.Uid, user *ptttype.UserecRaw, err error) {
	return cmbbs.PasswdLoadUser(userID)
}

func passwdSyncUpdate(uid ptttype.Uid, user *ptttype.UserecRaw) error {
	if uid < 1 || uid > ptttype.MAX_USERS {
		return cache.ErrInvalidUID
	}

	user.Money = cache.MoneyOf(uid)

	err := cmbbs.PasswdUpdate(uid, user)
	if err != nil {
		return err
	}

	return nil
}

func passwdSyncQuery(uid ptttype.Uid) (*ptttype.UserecRaw, error) {
	user, err := cmbbs.PasswdQuery(uid)
	if err != nil {
		return nil, err
	}

	user.Money = cache.MoneyOf(uid)

	return user, nil
}
