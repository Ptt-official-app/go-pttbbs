package ptt

import (
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

func pwcuStart(uid ptttype.Uid, userID *ptttype.UserID_t) (user *ptttype.UserecRaw, err error) {
	user, err = passwdSyncQuery(uid)
	if err != nil {
		return nil, err
	}

	if types.Cstrcmp(userID[:], user.UserID[:]) != 0 {
		return nil, ptttype.ErrInvalidUserID
	}

	return user, nil
}

func pwcuIncNumPost(user *ptttype.UserecRaw, uid ptttype.Uid) (err error) {
	u, err := pwcuStart(uid, &user.UserID)
	if err != nil {
		return err
	}

	u.NumPosts++
	user.NumPosts = u.NumPosts

	return pwcuEnd(uid, u)
}

func pwcuEnd(uid ptttype.Uid, user *ptttype.UserecRaw) (err error) {
	return passwdSyncUpdate(uid, user)
}

func pwcuRegCompleteJustify(uid ptttype.Uid, userID *ptttype.UserID_t, justify *ptttype.Reg_t) (err error) {
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

func pwcuBitEnableLevel(uid ptttype.Uid, userID *ptttype.UserID_t, perm ptttype.PERM) (err error) {
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

func pwcuBitDisableLevel(uid ptttype.Uid, userID *ptttype.UserID_t, perm ptttype.PERM) (err error) {
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
