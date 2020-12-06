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
