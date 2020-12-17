package bbs

import (
	"strconv"
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type UUserID string //The integrated bid-userid, concat with _, safe because bid is number >= 1.

func ToUUserID(uid ptttype.Uid, userIDRaw *ptttype.UserID_t) UUserID {
	if !uid.IsValid() {
		return UUserID("")
	}
	if !userIDRaw.IsValid() {
		return UUserID("")
	}
	return UUserID(uid.String() + "_" + types.CstrToString(userIDRaw[:]))
}

func (u UUserID) ToRaw() (uid ptttype.Uid, userIDRaw *ptttype.UserID_t, err error) {

	theList := strings.Split(string(u), "_")
	if len(theList) != 2 {
		return 0, nil, ErrInvalidUUserID
	}

	uid_i, err := strconv.Atoi(theList[0])
	if err != nil {
		return 0, nil, ErrInvalidUUserID
	}
	uid = ptttype.Uid(uid_i)
	if !uid.IsValid() {
		return 0, nil, ErrInvalidUUserID
	}

	userIDRaw = &ptttype.UserID_t{}
	copy(userIDRaw[:], []byte(theList[1]))
	if !userIDRaw.IsValid() {
		return 0, nil, ErrInvalidUUserID
	}
	return uid, userIDRaw, nil
}
