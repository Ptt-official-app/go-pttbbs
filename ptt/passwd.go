package ptt

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func InitCurrentUser(userID *ptttype.UserID_t) (uid ptttype.UID, user *ptttype.UserecRaw, err error) {
	uid, user, err = cmbbs.PasswdLoadUser(userID)
	if err != nil {
		return uid, user, err
	}

	// https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L736
	if types.Cstrcmp(user.UserID[:], []byte(ptttype.STR_GUEST)) == 0 {
		pwcuInitGuestPerm(user)
	}

	// https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L762
	if types.Cstrcmp(user.UserID[:], []byte(ptttype.STR_SYSOP)) == 0 {
		pwcuInitAdminPerm(user)
	}

	return uid, user, nil
}

func InitCurrentUserByUID(uid ptttype.UID) (user *ptttype.UserecRaw, err error) {
	user, err = cmbbs.PasswdQuery(uid)
	if err != nil {
		return nil, err
	}

	// https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L736
	if types.Cstrcmp(user.UserID[:], []byte(ptttype.STR_GUEST)) == 0 {
		pwcuInitGuestPerm(user)
	}

	// https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L762
	if types.Cstrcmp(user.UserID[:], ptttype.STR_SYSOP) == 0 {
		pwcuInitAdminPerm(user)
	}

	return user, nil
}

func passwdSyncUpdate(uid ptttype.UID, user *ptttype.UserecRaw) error {
	if !uid.IsValid() {
		return cache.ErrInvalidUID
	}

	user.Money = cache.MoneyOf(uid)

	err := cmbbs.PasswdUpdate(uid, user)
	if err != nil {
		return err
	}

	return nil
}

func passwdSyncQuery(uid ptttype.UID) (*ptttype.UserecRaw, error) {
	user, err := cmbbs.PasswdQuery(uid)
	if err != nil {
		return nil, err
	}

	user.Money = cache.MoneyOf(uid)

	return user, nil
}

// pwcuLoginSave update user numLoginDays, LastLogin and LastSeen
func pwcuLoginSave(uid ptttype.UID, user *ptttype.UserecRaw, ip *ptttype.IPv4_t) (isFirstLoginOfDay bool, err error) {
	// get user 1st login
	firstLoginDay := user.FirstLogin.ToLocal()
	// get 1st day at 00:00
	baseRefTime := types.TimeToTime4(time.Date(firstLoginDay.Year(), firstLoginDay.Month(), firstLoginDay.Day(), 0, 0, 0, 0, firstLoginDay.Location()))
	loginStartTime := types.NowTS()
	refTime := loginStartTime

	// multiple login?
	if refTime < user.LastLogin {
		refTime = user.LastLogin
	}

	regDays := (refTime - baseRefTime) / 86400
	prevRegDays := (user.LastLogin - baseRefTime) / 86400
	// error check?
	if uint32(user.NumLoginDays) > uint32(prevRegDays)+1 {
		user.NumLoginDays = uint32(prevRegDays) + 1
	}
	logrus.Info("refTime: ", refTime)
	logrus.Info("baseRefTime: ", baseRefTime)
	logrus.Info("user.LastLogin: ", user.LastLogin)
	logrus.Info("regDays: ", regDays)
	logrus.Info("prevRegDays: ", prevRegDays)
	if regDays > prevRegDays {
		user.NumLoginDays++
		isFirstLoginOfDay = true
	}

	user.LastLogin = loginStartTime
	user.LastSeen = loginStartTime
	user.LastHost = *ip
	err = passwdSyncUpdate(uid, user)
	if err != nil {
		return isFirstLoginOfDay, err
	}

	return isFirstLoginOfDay, nil
}
