package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func LoadGeneralArticles(uuserID UUserID, bboardID BBoardID, startIdxStr string, nArticles int) (summary []*ArticleSummary, nextIdxStr string, isNewest bool, err error) {

	if nArticles < 1 {
		return nil, "", false, ErrInvalidParams
	}

	startIdx, err := ptttype.ToSortIdx(startIdxStr)
	if err != nil {
		return nil, "", false, ErrInvalidParams
	}
	if startIdx < 0 {
		return nil, "", false, ErrInvalidParams
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, "", false, ErrInvalidParams
	}

	uid, _, err := uuserID.ToRaw()
	if err != nil {
		return nil, "", false, ErrInvalidParams
	}

	userecRaw, err := ptt.InitCurrentUserByUid(uid)
	if err != nil {
		return nil, "", false, err
	}

	summaryRaw, nextIdx, isNewest, err := ptt.LoadGeneralArticles(userecRaw, uid, boardIDRaw, bid, startIdx, nArticles)
	if err != nil {
		return nil, "", false, err
	}

	summary = make([]*ArticleSummary, len(summaryRaw))
	for idx, each := range summaryRaw {
		eachSummary := NewArticleSummaryFromRaw(bboardID, each)
		summary[idx] = eachSummary
	}

	nextIdxStr = nextIdx.String()

	return summary, nextIdxStr, isNewest, nil
}
