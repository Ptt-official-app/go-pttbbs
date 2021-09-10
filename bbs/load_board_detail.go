package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func LoadBoardDetail(uuserID UUserID, bboardID BBoardID) (detail *BoardDetail, err error) {
	bid, _, err := bboardID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}
	userID, err := uuserID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}
	uid, userecRaw, err := ptt.InitCurrentUser(userID)
	if err != nil {
		return nil, err
	}

	boardDetailRaw, err := ptt.LoadBoardDetail(userecRaw, uid, bid)
	if err != nil {
		return nil, err
	}

	detail = NewBoardDetailFromRaw(boardDetailRaw, bid)

	return detail, nil
}
