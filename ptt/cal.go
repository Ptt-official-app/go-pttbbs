package ptt

import "github.com/Ptt-official-app/go-pttbbs/ptttype"

func getRestrictionReason(numLoginDays uint32, badPost uint8, postLimitLogins uint8, postLimitBadpost uint8) (reason ptttype.RestrictReason, err error) {

	if numLoginDays/10 < uint32(postLimitLogins) {
		return ptttype.RESTRICT_REASON_NUMLOGIN_DAYS, nil
	}

	if badPost > (255 - postLimitBadpost) {
		return ptttype.RESTRICT_REASON_BADPOST, nil
	}

	return ptttype.RESTRICT_REASON_NONE, nil
}
