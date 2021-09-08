package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// LoadGeneralArticles in descending mode.
func LoadGeneralArticles(
	uuserID UUserID,
	bboardID BBoardID,
	startIdxStr string,
	nArticles int,
	isDesc bool) (

	summaries []*ArticleSummary,
	nextIdxStr string,
	nextCreateTime types.Time4,
	isNewest bool,
	startNumIdx ptttype.SortIdx,
	err error) {

	if nArticles < 1 {
		return nil, "", 0, false, -1, ErrInvalidParams
	}

	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, "", 0, false, -1, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, "", 0, false, -1, err
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, "", 0, false, -1, ErrInvalidParams
	}

	// 1. find start idx. start-idx as nextCreateTime if unable to find startIdxStr
	startIdx, err := loadGeneralArticlesToStartIdx(userecRaw, uid, boardIDRaw, bid, startIdxStr, isDesc)
	if err != nil {
		return nil, "", 0, false, -1, err
	}

	// 2. load articles.
	var summariesRaw []*ptttype.ArticleSummaryRaw
	var nextSummaryRaw *ptttype.ArticleSummaryRaw
	if startIdx >= 0 {
		summariesRaw, isNewest, nextSummaryRaw, startNumIdx, err = ptt.LoadGeneralArticles(userecRaw, uid, boardIDRaw, bid, startIdx, nArticles, isDesc)
		if err != nil {
			return nil, "", 0, false, -1, err
		}
	}

	// 3. nextIdxStr
	nextIdxStr = ""
	nextCreateTime = 0
	if nextSummaryRaw != nil {
		nextSummary := NewArticleSummaryFromRaw(bboardID, nextSummaryRaw)
		nextIdxStr = nextSummary.Idx
		nextCreateTime = nextSummary.CreateTime
	}

	summaries = make([]*ArticleSummary, len(summariesRaw))
	for idx, each := range summariesRaw {
		eachSummary := NewArticleSummaryFromRaw(bboardID, each)
		summaries[idx] = eachSummary
	}

	return summaries, nextIdxStr, nextCreateTime, isNewest, startNumIdx, nil
}

func loadGeneralArticlesToStartIdx(
	userecRaw *ptttype.UserecRaw,
	uid ptttype.UID,
	boardIDRaw *ptttype.BoardID_t,
	bid ptttype.Bid,
	startIdxStr string,
	isDesc bool,
) (startIdx ptttype.SortIdx, err error) {
	if startIdxStr == "" {
		if isDesc {
			return 0, nil
		} else {
			return 1, nil
		}
	}

	createTime, articleID, err := DeserializeArticleIdxStr(startIdxStr)
	if err != nil {
		return -1, err
	}
	filename := articleID.ToRaw()

	startIdx, err = ptt.FindArticleStartIdx(userecRaw, uid, boardIDRaw, bid, createTime, filename, isDesc)

	return startIdx, err
}
