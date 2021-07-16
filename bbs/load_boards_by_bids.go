package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func LoadBoardsByBids(uuserID UUserID, bids []ptttype.Bid) (summaries []*BoardSummary, err error) {
	userID, err := uuserID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userID)
	if err != nil {
		return nil, err
	}

	summaryRaw, err := ptt.LoadBoardsByBids(userecRaw, uid, bids)
	if err != nil {
		return nil, err
	}

	summaries = make([]*BoardSummary, len(summaryRaw))
	for idx, each := range summaryRaw {
		eachSummary := NewBoardSummaryFromRaw(each)
		summaries[idx] = eachSummary
	}

	return summaries, nil
}
