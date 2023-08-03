package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

// valid users to see/change the user-info.
func userInfoIsValidUser(uuserID bbs.UUserID, queryUUserID bbs.UUserID) (isValid bool) {
	if queryUUserID == bbs.UUserID(ptttype.STR_GUEST) {
		return false
	}

	permSysop := ptttype.PERM_ACCOUNTS | ptttype.PERM_SYSOP | ptttype.PERM_ACCTREG
	if bbs.IsSysop(uuserID, permSysop) {
		return true
	}

	return uuserID == queryUUserID
}

// valid users to see/change email / user-level2
func userInfoIsValidEmailUser(uuserID bbs.UUserID, queryUUserID bbs.UUserID, jwt string, context EmailTokenContext, isAllowSysop bool) (isValid bool, email string) {
	if queryUUserID == bbs.UUserID(ptttype.STR_GUEST) {
		return false, ""
	}

	isSysop := false
	if isAllowSysop {
		permSysop := ptttype.PERM_ACCOUNTS | ptttype.PERM_SYSOP | ptttype.PERM_ACCTREG
		isSysop = bbs.IsSysop(uuserID, permSysop)
	}

	if !isSysop && uuserID != queryUUserID {
		return false, ""
	}

	emailUserID, _, _, email, err := VerifyEmailJwt(jwt, context)
	if err != nil {
		return false, ""
	}

	if queryUUserID != emailUserID {
		return false, ""
	}

	return true, email
}
