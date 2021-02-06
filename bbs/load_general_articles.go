package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//LoadGeneralArticles in descending mode.
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

	//1. find start idx. start-idx as nextCreateTime if unable to find startIdxStr
	startIdx, err := loadGeneralArticlesToStartIdx(userecRaw, uid, boardIDRaw, bid, startIdxStr, isDesc)
	if err != nil {
		return nil, "", 0, false, -1, ErrInvalidParams
	}

	//2. load articles.
	var summariesRaw []*ptttype.ArticleSummaryRaw
	var nextSummaryRaw *ptttype.ArticleSummaryRaw
	if startIdx >= 0 {
		summariesRaw, isNewest, nextSummaryRaw, startNumIdx, err = ptt.LoadGeneralArticles(userecRaw, uid, boardIDRaw, bid, startIdx, nArticles, isDesc)
		if err != nil {
			return nil, "", 0, false, -1, err
		}
	}

	//3. nextIdxStr
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

	if startIdxStr == "" { //no need to check startIdxStr
		return summaries, nextIdxStr, nextCreateTime, isNewest, startNumIdx, nil
	}

	// check same-create-time
	if len(summaries) > 0 && summaries[0].Idx == startIdxStr {
		return summaries, nextIdxStr, nextCreateTime, isNewest, startNumIdx, nil
	}

	summariesSameCreateTime, startNumIdxSameCreateTime, err := loadGeneralArticlesSameCreateTime(userecRaw, uid, bboardID, boardIDRaw, bid, startIdxStr, isDesc)
	if err != nil || len(summariesSameCreateTime) == 0 {
		return summaries, nextIdxStr, nextCreateTime, isNewest, startNumIdx, nil
	}

	summaries = append(summariesSameCreateTime, summaries...)
	startNumIdx = startNumIdxSameCreateTime

	return summaries, nextIdxStr, nextCreateTime, isNewest, startNumIdx, nil
}

func loadGeneralArticlesToStartIdx(
	userecRaw *ptttype.UserecRaw,
	uid ptttype.Uid,
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

	createTime, articleID, err := deserializeArticleIdxStr(startIdxStr)
	if err != nil {
		return -1, err
	}
	filename, _ := articleID.ToRaw()

	startIdx, err = ptt.FindArticleStartIdx(userecRaw, uid, boardIDRaw, bid, createTime, filename, isDesc)

	return startIdx, err
}

func loadGeneralArticlesSameCreateTime(userecRaw *ptttype.UserecRaw, uid ptttype.Uid, boardID BBoardID, boardIDRaw *ptttype.BoardID_t, bid ptttype.Bid, startIdxStr string, isDesc bool) (summaries []*ArticleSummary, startNumIdx ptttype.SortIdx, err error) {
	createTime, _, err := deserializeArticleIdxStr(startIdxStr)
	if err != nil {
		return nil, 0, err
	}

	startAid, err := ptt.FindArticleStartIdx(userecRaw, uid, boardIDRaw, bid, createTime, nil, true)
	if err != nil {
		return nil, 0, err
	}
	if startAid == -1 {
		startAid = 1
	}
	endAid, err := ptt.FindArticleStartIdx(userecRaw, uid, boardIDRaw, bid, createTime, nil, false)
	if err != nil {
		return nil, 0, err
	}
	if endAid == -1 {
		endAid = 0
	}

	summariesRaw, newStartNumIdx, newEndNumIdx, err := ptt.LoadGeneralArticlesSameCreateTime(boardIDRaw, bid, startAid, endAid, createTime)
	if err != nil {
		return nil, 0, err
	}

	if len(summariesRaw) == 0 {
		return nil, 0, nil
	}

	summaries = make([]*ArticleSummary, len(summariesRaw))
	for idx, each := range summariesRaw {
		summaries[idx] = NewArticleSummaryFromRaw(boardID, each)
	}

	if isDesc {
		reverseArticleSummaries(summaries)
		startNumIdx = newEndNumIdx
	} else {
		startNumIdx = newStartNumIdx
	}

	return summaries, startNumIdx, nil
}

func reverseArticleSummaries(summaries []*ArticleSummary) {
	for i, j := 0, len(summaries)-1; i < j; i, j = i+1, j-1 {
		summaries[i], summaries[j] = summaries[j], summaries[i]
	}
}
