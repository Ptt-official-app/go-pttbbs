package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func LoadGeneralBoardDetails(uuserID UUserID, startIdxStr string, nBoards int, isAsc bool, bsortBy ptttype.BSortBy) (details []*BoardDetail, nextIdxStr string, err error) {
	startIdx, err := loadGeneralBoardsToStartIdx(startIdxStr, isAsc, bsortBy)
	if err != nil {
		return nil, "", ErrInvalidParams
	}
	if startIdx < 0 {
		return nil, "", nil
	}

	userID, err := uuserID.ToRaw()
	if err != nil {
		return nil, "", ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userID)
	if err != nil {
		return nil, "", err
	}

	detailsRaw, nextDetailRaw, err := ptt.LoadGeneralBoardDetails(userecRaw, uid, startIdx, nBoards, isAsc, bsortBy)
	if err != nil {
		return nil, "", err
	}

	details = make([]*BoardDetail, len(detailsRaw))
	for idx, each := range detailsRaw {
		eachDetail := NewBoardDetailFromRaw(each, each.Bid)
		details[idx] = eachDetail
	}
	if nextDetailRaw != nil {
		nextSummary := NewBoardDetailFromRaw(nextDetailRaw, nextDetailRaw.Bid)
		switch bsortBy {
		case ptttype.BSORT_BY_NAME:
			nextIdxStr = nextSummary.IdxByName
		case ptttype.BSORT_BY_CLASS:
			nextIdxStr = nextSummary.IdxByClass
		}
	}

	return details, nextIdxStr, nil
}
