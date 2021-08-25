package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func LoadClassBoards(uuserID UUserID, clsBid ptttype.Bid, bsortBy ptttype.BSortBy) (summaries []*BoardSummary, err error) {
	userID, err := uuserID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userID)
	if err != nil {
		return nil, err
	}

	summariesRaw, err := ptt.LoadClassBoards(userecRaw, uid, clsBid, bsortBy)
	if err != nil {
		return nil, err
	}

	summaries = make([]*BoardSummary, len(summariesRaw))
	for idx, each := range summariesRaw {
		eachSummary := NewBoardSummaryFromRaw(each)
		summaries[idx] = eachSummary
	}

	return summaries, nil
}
