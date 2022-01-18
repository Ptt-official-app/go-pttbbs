package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/boardd"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func LoadHotBoards(uuserID UUserID) (summary []*BoardSummary, err error) {
	userID, err := uuserID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userID)
	if err != nil {
		return nil, err
	}

	summaryRaw, err := boardd.LoadHotBoards(userecRaw, uid)
	if err != nil {
		return nil, err
	}

	summary = make([]*BoardSummary, len(summaryRaw))
	for idx, each := range summaryRaw {
		eachSummary := NewBoardSummaryFromRaw(each)
		summary[idx] = eachSummary
	}

	return summary, nil
}
