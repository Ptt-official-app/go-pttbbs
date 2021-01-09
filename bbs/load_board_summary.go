package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func LoadBoardSummary(uuserID UUserID, bboardID BBoardID) (summary *BoardSummary, err error) {
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

	summaryRaw, err := ptt.LoadBoardSummary(userecRaw, uid, bid)
	if err != nil || summaryRaw == nil {
		return nil, err
	}
	summary = NewBoardSummaryFromRaw(summaryRaw)

	return summary, nil
}
