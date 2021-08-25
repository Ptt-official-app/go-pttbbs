package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func LoadFullClassBoards(uuserID UUserID, startBid ptttype.Bid, nBoards int) (summaries []*BoardSummary, nextBid ptttype.Bid, err error) {
	if !startBid.IsValid() {
		return nil, 0, ptttype.ErrInvalidBid
	}

	userID, err := uuserID.ToRaw()
	if err != nil {
		return nil, 0, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userID)
	if err != nil {
		return nil, 0, err
	}

	summariesRaw, nextSummaryRaw, err := ptt.LoadFullClassBoards(userecRaw, uid, startBid, nBoards)
	if err != nil {
		return nil, 0, err
	}

	summaries = make([]*BoardSummary, len(summariesRaw))
	for idx, each := range summariesRaw {
		eachSummary := NewBoardSummaryFromRaw(each)
		summaries[idx] = eachSummary
	}

	if nextSummaryRaw != nil {
		nextBid = nextSummaryRaw.Bid
	}

	return summaries, nextBid, nil
}
