package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func LoadGeneralBoards(uuserID UUserID, startIdxStr string, nBoards int, title []byte, keyword []byte, isAsc bool, bsortBy ptttype.BSortBy) (summaries []*BoardSummary, nextIdxStr string, err error) {
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

	summariesRaw, nextSummaryRaw, err := ptt.LoadGeneralBoards(userecRaw, uid, startIdx, nBoards, title, keyword, isAsc, bsortBy)
	if err != nil {
		return nil, "", err
	}

	summaries = make([]*BoardSummary, len(summariesRaw))
	for idx, each := range summariesRaw {
		eachSummary := NewBoardSummaryFromRaw(each)
		summaries[idx] = eachSummary
	}
	if nextSummaryRaw != nil {
		nextSummary := NewBoardSummaryFromRaw(nextSummaryRaw)
		switch bsortBy {
		case ptttype.BSORT_BY_NAME:
			nextIdxStr = nextSummary.IdxByName
		case ptttype.BSORT_BY_CLASS:
			nextIdxStr = nextSummary.IdxByClass
		}
	}

	return summaries, nextIdxStr, nil
}

func loadGeneralBoardsToStartIdx(startIdxStr string, isAsc bool, bsortBy ptttype.BSortBy) (startIdx ptttype.SortIdx, err error) {
	if startIdxStr == "" {
		if isAsc {
			return 1, nil
		} else {
			return 0, nil
		}
	}

	switch bsortBy {
	case ptttype.BSORT_BY_NAME:
		brdname, err := deserializeBoardIdxByNameStr(startIdxStr)
		if err != nil {
			return -1, err
		}
		brdnameRaw := &ptttype.BoardID_t{}
		copy(brdnameRaw[:], []byte(brdname))
		startIdx, err = ptt.FindBoardStartIdxByName(brdnameRaw, isAsc)
		if err != nil {
			return -1, err
		}
		return startIdx, nil
	case ptttype.BSORT_BY_CLASS:
		cls, brdname, err := deserializeBoardIdxByClassStr(startIdxStr)
		if err != nil {
			return -1, err
		}
		brdnameRaw := &ptttype.BoardID_t{}
		copy(brdnameRaw[:], []byte(brdname))
		startIdx, err = ptt.FindBoardStartIdxByClass(cls, brdnameRaw, isAsc)
		if err != nil {
			return -1, err
		}
		return startIdx, nil
	default:
		return -1, ErrInvalidParams
	}
}
