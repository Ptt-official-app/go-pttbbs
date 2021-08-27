package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func IsBoardsValidUser(uuserID UUserID, boardIDs []BBoardID) (isValid map[BBoardID]bool, err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, err
	}

	isValid = make(map[BBoardID]bool)

	for _, eachBoardID := range boardIDs {
		bid, boardIDRaw, err := eachBoardID.ToRaw()
		if err != nil {
			continue
		}
		eachIsValid, err := ptt.IsBoardValidUser(userecRaw, uid, boardIDRaw, bid)
		if err != nil {
			continue
		}
		isValid[eachBoardID] = eachIsValid
	}

	return isValid, nil
}
