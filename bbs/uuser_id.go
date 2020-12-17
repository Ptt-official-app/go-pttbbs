package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type UUserID string //The integrated bid-userid, concat with _, safe because bid is number >= 1.

func ToUUserID(userIDRaw *ptttype.UserID_t) UUserID {
	if !userIDRaw.IsValid() {
		return UUserID("")
	}
	return UUserID(types.CstrToString(userIDRaw[:]))
}

func (u UUserID) ToRaw() (userIDRaw *ptttype.UserID_t, err error) {
	userIDRaw = &ptttype.UserID_t{}
	copy(userIDRaw[:], []byte(u))
	if !userIDRaw.IsValid() {
		return nil, ErrInvalidUUserID
	}
	return userIDRaw, nil
}
