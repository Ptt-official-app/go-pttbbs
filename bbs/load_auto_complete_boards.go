package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func LoadAutoCompleteBoards(uuserID UUserID, startIdxStr string, nBoards int, keyword string, isAsc bool) (summaries []*BoardSummary, nextIdxStr string, err error) {

	keywordBytes := []byte(keyword)

	startIdx, err := loadAutoCompleteBoardsToStartIdx(startIdxStr, keywordBytes, isAsc)
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

	summariesRaw, nextSummaryRaw, err := ptt.LoadAutoCompleteBoards(userecRaw, uid, startIdx, nBoards, keywordBytes, isAsc)
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
		nextIdxStr = nextSummary.IdxByName
	}

	return summaries, nextIdxStr, nil
}

func loadAutoCompleteBoardsToStartIdx(startIdxStr string, keywordBytes []byte, isAsc bool) (startIdx ptttype.SortIdx, err error) {
	if startIdxStr != "" {
		return loadGeneralBoardsToStartIdx(startIdxStr, isAsc, ptttype.BSORT_BY_NAME)
	}

	if len(keywordBytes) == 0 {
		return 1, nil
	}

	return ptt.FindBoardAutoCompleteStartIdx(keywordBytes, isAsc)
}
