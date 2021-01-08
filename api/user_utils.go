package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

//valid users to see/change the user-info.
func userIsValidUser(uuserID bbs.UUserID, queryUUserID bbs.UUserID) (isValid bool) {
	permSysop := ptttype.PERM_ACCOUNTS | ptttype.PERM_SYSOP | ptttype.PERM_ACCTREG
	if bbs.IsSysop(uuserID, permSysop) {
		return true
	}

	return uuserID == queryUUserID
}

//valid users to see/change email / user-level2
func userIsValidEmailUser(uuserID bbs.UUserID, queryUUserID bbs.UUserID, jwt string, isSysop bool) (isValid bool, email string) {

	if isSysop {
		permSysop := ptttype.PERM_ACCOUNTS | ptttype.PERM_SYSOP | ptttype.PERM_ACCTREG
		if bbs.IsSysop(uuserID, permSysop) {
			return true, ""
		}
	}

	if uuserID != queryUUserID {
		return false, ""
	}

	emailUserID, _, email, err := VerifyEmailJwt(jwt)
	if err != nil {
		return false, ""
	}

	if queryUUserID != emailUserID {
		return false, ""
	}

	return true, email
}
