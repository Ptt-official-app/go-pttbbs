package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func LoadGeneralArticles(userID string, bboardID BBoardID, startIdxStr string, nArticles int) (summary []*ArticleSummary, nextIdxStr string, err error) {

	if nArticles < 1 {
		return nil, "", ErrInvalidParams
	}

	startIdx, err := ptttype.ToSortIdx(startIdxStr)
	if err != nil {
		return nil, "", ErrInvalidParams
	}
	if startIdx < 0 {
		return nil, "", ErrInvalidParams
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, "", ErrInvalidParams
	}

	userIDRaw := &ptttype.UserID_t{}
	copy(userIDRaw[:], []byte(userID))

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, "", err
	}

	summaryRaw, nextIdx, err := ptt.LoadGeneralArticles(userecRaw, uid, boardIDRaw, bid, startIdx, nArticles)
	if err != nil {
		return nil, "", err
	}

	summary = make([]*ArticleSummary, len(summaryRaw))
	for idx, each := range summaryRaw {
		eachSummary := NewArticleSummaryFromRaw(each)
		summary[idx] = eachSummary
	}

	nextIdxStr = nextIdx.String()

	return summary, nextIdxStr, nil
}
