package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func LoadGeneralBoards(uuserID UUserID, startIdxStr string, nBoards int, keyword string) (summary []*BoardSummary, nextIdxStr string, err error) {
	startIdx, err := ptttype.ToSortIdx(startIdxStr)
	if err != nil {
		return nil, "", ErrInvalidParams
	}
	if startIdx < 0 {
		return nil, "", ErrInvalidParams
	}

	uid, _, err := uuserID.ToRaw()
	if err != nil {
		return nil, "", ErrInvalidParams
	}

	userecRaw, err := ptt.InitCurrentUserByUid(uid)
	if err != nil {
		return nil, "", err
	}

	keywordBytes := types.Utf8ToBig5(keyword)

	summaryRaw, nextIdx, err := ptt.LoadGeneralBoards(userecRaw, uid, startIdx, nBoards, keywordBytes)
	if err != nil {
		return nil, "", err
	}

	summary = make([]*BoardSummary, len(summaryRaw))
	for idx, each := range summaryRaw {
		eachSummary := NewBoardSummaryFromRaw(each)
		summary[idx] = eachSummary
	}

	nextIdxStr = nextIdx.String()

	return summary, nextIdxStr, nil
}
