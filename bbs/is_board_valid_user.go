package bbs

import "github.com/Ptt-official-app/go-pttbbs/ptt"

func IsBoardValidUser(uuserID UUserID, boardID BBoardID) (isValid bool, err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return false, ErrInvalidParams
	}

	bid, boardIDRaw, err := boardID.ToRaw()
	if err != nil {
		return false, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return false, err
	}

	return ptt.IsBoardValidUser(userecRaw, uid, boardIDRaw, bid)
}
