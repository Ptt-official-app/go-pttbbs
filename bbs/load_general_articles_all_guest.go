package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// LoadGeneralArticles in descending mode.
func LoadGeneralArticlesAllGuest(
	bboardID BBoardID,
	startIdxStr string,
	nArticles int,
	isDesc bool) (
	summaries []*ArticleSummary,
	nextIdxStr string,
	nextCreateTime types.Time4,
	isNewest bool,
	startNumIdx ptttype.SortIdx,
	err error,
) {
	if nArticles < 1 {
		return nil, "", 0, false, -1, ErrInvalidParams
	}

	boardIDRaw, err := bboardID.ToRawAllGuest()
	if err != nil {
		return nil, "", 0, false, -1, err
	}

	// 1. find start idx. start-idx as nextCreateTime if unable to find startIdxStr
	startIdx, err := loadGeneralArticlesToStartIdxAllGuest(boardIDRaw, startIdxStr, isDesc)
	if err != nil {
		return nil, "", 0, false, -1, err
	}

	// 2. load articles.
	var summariesRaw []*ptttype.ArticleSummaryRaw
	var nextSummaryRaw *ptttype.ArticleSummaryRaw
	if startIdx >= 0 {
		summariesRaw, isNewest, nextSummaryRaw, startNumIdx, err = ptt.LoadGeneralArticlesAllGuest(boardIDRaw, startIdx, nArticles, isDesc)
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

func loadGeneralArticlesToStartIdxAllGuest(
	boardIDRaw *ptttype.BoardID_t,
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

	startIdx, err = ptt.FindArticleStartIdxAllGuest(boardIDRaw, createTime, filename, isDesc)

	return startIdx, err
}
