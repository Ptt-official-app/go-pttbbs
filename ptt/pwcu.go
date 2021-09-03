package ptt

import (
	"time"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func pwcuEnableBit(perm ptttype.PERM, mask ptttype.PERM) ptttype.PERM {
	return perm | mask
}

func pwcuDisableBit(perm ptttype.PERM, mask ptttype.PERM) ptttype.PERM {
	return perm & ^mask
}

func pwcuEnableUFlagBit(uflag ptttype.UFlag, mask ptttype.UFlag) ptttype.UFlag {
	return uflag | mask
}

func pwcuSetByBit(perm ptttype.PERM, mask ptttype.PERM, isSet bool) ptttype.PERM {
	if isSet {
		return pwcuEnableBit(perm, mask)
	} else {
		return pwcuDisableBit(perm, mask)
	}
}

func pwcuStart(uid ptttype.UID, userID *ptttype.UserID_t) (user *ptttype.UserecRaw, err error) {
	user, err = passwdSyncQuery(uid)
	if err != nil {
		return nil, err
	}

	if types.Cstrcmp(userID[:], user.UserID[:]) != 0 {
		return nil, ptttype.ErrInvalidUserID
	}

	return user, nil
}

func pwcuIncNumPost(user *ptttype.UserecRaw, uid ptttype.UID) (err error) {
	u, err := pwcuStart(uid, &user.UserID)
	if err != nil {
		return err
	}

	u.NumPosts++
	user.NumPosts = u.NumPosts

	return pwcuEnd(uid, u)
}

func pwcuEnd(uid ptttype.UID, user *ptttype.UserecRaw) (err error) {
	return passwdSyncUpdate(uid, user)
}

func pwcuRegCompleteJustify(uid ptttype.UID, userID *ptttype.UserID_t, justify *ptttype.Reg_t) (err error) {
	var user *ptttype.UserecRaw

	user, err = pwcuStart(uid, userID)
	if err != nil {
		return err
	}
	defer func() {
		err = pwcuEnd(uid, user)
	}()

	copy(user.Justify[:], justify[:])
	user.UserLevel = pwcuEnableBit(user.UserLevel, ptttype.PERM_POST|ptttype.PERM_LOGINOK)

	return nil
}

func pwcuBitEnableLevel(uid ptttype.UID, userID *ptttype.UserID_t, perm ptttype.PERM) (err error) {
	var user *ptttype.UserecRaw
	user, err = pwcuStart(uid, userID)
	if err != nil {
		return err
	}
	defer func() {
		errEnd := pwcuEnd(uid, user)
		if err == nil {
			err = errEnd
		}
	}()

	_ = pwcuEnableBit(user.UserLevel, perm)

	return nil
}

func pwcuBitDisableLevel(uid ptttype.UID, userID *ptttype.UserID_t, perm ptttype.PERM) (err error) {
	var user *ptttype.UserecRaw

	user, err = pwcuStart(uid, userID)
	if err != nil {
		return err
	}
	defer func() {
		err = pwcuEnd(uid, user)
	}()

	_ = pwcuDisableBit(user.UserLevel, perm)

	return nil
}

func pwcuInitGuestPerm(user *ptttype.UserecRaw) {
	user.UserLevel = 0
	user.UFlag = ptttype.UF_BRDSORT
	user.Pager = ptttype.PAGER_OFF
	if ptttype.DBCSAWARE {
		user.UFlag = pwcuEnableUFlagBit(user.UFlag, ptttype.UF_DBCS_AWARE|ptttype.UF_DBCS_DROP_REPEAT)
		if ptttype.GUEST_DEFAULT_DBCS_NOINTRESC {
			user.UFlag = pwcuEnableUFlagBit(user.UFlag, ptttype.UF_DBCS_NOINTRESC)
		}
	}
}

func pwcuInitAdminPerm(user *ptttype.UserecRaw) {
	user.UserLevel = ptttype.PERM_BASIC | ptttype.PERM_CHAT | ptttype.PERM_PAGE |
		ptttype.PERM_POST | ptttype.PERM_LOGINOK | ptttype.PERM_MAILLIMIT |
		ptttype.PERM_CLOAK | ptttype.PERM_SEECLOAK | ptttype.PERM_XEMPT |
		ptttype.PERM_SYSOPHIDE | ptttype.PERM_BM | ptttype.PERM_ACCOUNTS |
		ptttype.PERM_CHATROOM | ptttype.PERM_BOARD | ptttype.PERM_SYSOP | ptttype.PERM_BBSADM
}

// pwcuLoginSave
// update user numLoginDays, LastLogin and LastSeen
// https://github.com/ptt/pttbbs/blob/master/mbbsd/passwd.c#L501
func pwcuLoginSave(uid ptttype.UID, user *ptttype.UserecRaw, uinfo *ptttype.UserInfoRaw, ip *ptttype.IPv4_t) (isFirstLoginOfDay bool, err error) {
	var u *ptttype.UserecRaw

	u, err = pwcuStart(uid, &user.UserID)
	if err != nil {
		return false, err
	}
	defer func() {
		err = pwcuEnd(uid, u)
	}()

	// new host from 'fromhost'
	// https://github.com/ptt/pttbbs/blob/master/mbbsd/passwd.c#L517
	u.LastHost = *ip
	user.LastHost = *ip

	// get user 1st login
	firstLoginDay := u.FirstLogin.ToLocal()
	// get 1st day at 00:00
	baseRefTime := types.TimeToTime4(time.Date(firstLoginDay.Year(), firstLoginDay.Month(), firstLoginDay.Day(), 0, 0, 0, 0, firstLoginDay.Location()))

	// loginStartTime
	loginStartTime := types.NowTS()
	refTime := loginStartTime

	// multiple login?
	if refTime < u.LastLogin {
		refTime = u.LastLogin
	}

	regDays := (refTime - baseRefTime) / ptttype.DAY_SECONDS
	prevRegDays := (u.LastLogin - baseRefTime) / ptttype.DAY_SECONDS
	// error check?
	// plus one for initial day
	// https://github.com/ptt/pttbbs/blob/master/mbbsd/passwd.c#L540
	if uint32(u.NumLoginDays) > uint32(prevRegDays)+1 {
		u.NumLoginDays = uint32(prevRegDays) + 1
	}

	// calculate numlogindays (only increase one per each key)
	// https://github.com/ptt/pttbbs/blob/master/mbbsd/passwd.c#L544
	if regDays > prevRegDays {
		u.NumLoginDays++
		isFirstLoginOfDay = true
	}

	user.NumLoginDays = u.NumLoginDays

	// update last login time
	// https://github.com/ptt/pttbbs/blob/master/mbbsd/passwd.c#L552
	u.LastLogin = loginStartTime
	user.LastLogin = loginStartTime

	if !uinfo.UserLevel.Hide() {
		u.LastSeen = loginStartTime
		user.LastSeen = loginStartTime
	}

	return isFirstLoginOfDay, nil
}
