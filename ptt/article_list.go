package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

//LoadGeneralArticles
//
//load general articles in descending mode.
//
//Implementing:
//1. preprocessing in i_read (perm-check, cache-preparation)
//2. get_records_and_bottom (with only records)
//https://github.com/ptt/pttbbs/blob/master/mbbsd/read.c#L1197
//https://github.com/ptt/pttbbs/blob/master/mbbsd/read.c#L1106
//
//get bottom can be a separated api.
func LoadGeneralArticles(user *ptttype.UserecRaw, uid ptttype.Uid, boardIDRaw *ptttype.BoardID_t, bid ptttype.Bid, startAid ptttype.Aid, nArticles int, isDesc bool) (summaries []*ptttype.ArticleSummaryRaw, isNewest bool, nextSummary *ptttype.ArticleSummaryRaw, startNumIdx ptttype.NumIdx, err error) {

	//1. check perm.
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, false, nil, 0, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, false, nil, 0, ErrNotPermitted
	}

	//2. bcache preparation.
	total, err := cache.GetBTotalWithRetry(bid)
	if err != nil {
		return nil, false, nil, 0, err
	}
	if total == 0 {
		return nil, true, nil, 0, nil
	}

	//3. get records
	filename, err := setBDir(boardIDRaw)
	if err != nil {
		return nil, false, nil, 0, err
	}
	//3.1. ensure startAid
	if startAid == 0 && isDesc {
		startAid = ptttype.Aid(total)
	}

	maxAid := ptttype.Aid(total)
	summaries, err = cmsys.GetRecords(boardIDRaw, filename, startAid, nArticles+1, isDesc, maxAid)
	if err != nil {
		return nil, false, nil, 0, err
	}

	//4. return
	isNewest = false
	if isDesc {
		isNewest = startAid == ptttype.Aid(total)
	} else {
		isNewest = len(summaries) != nArticles+1
	}

	if len(summaries) == nArticles+1 {
		nextSummary = summaries[len(summaries)-1]
		summaries = summaries[:len(summaries)-1]
	}

	logrus.Infof("ptt.LoadGeneralArticles: to return: startAid: %v", startAid)

	return summaries, isNewest, nextSummary, ptttype.NumIdx(startAid), nil
}

func FindArticleStartAid(user *ptttype.UserecRaw, uid ptttype.Uid, boardID *ptttype.BoardID_t, bid ptttype.Bid, createTime types.Time4, filename *ptttype.Filename_t, isDesc bool) (startAid ptttype.Aid, err error) {

	//1. check perm.
	board, err := cache.GetBCache(bid)
	if err != nil {
		return -1, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return -1, ErrNotPermitted
	}

	//3. get records
	dirFilename, err := setBDir(boardID)
	if err != nil {
		return -1, err
	}

	total, err := cache.GetBTotalWithRetry(bid)
	if err != nil {
		return -1, err
	}
	if total == 0 {
		return -1, nil
	}

	logrus.Infof("FindArticleStartAid: to cmsys.FindRecordStartAid: dirFilename: %v total: %v createTime: %v filename: %v", dirFilename, total, createTime, filename)

	return cmsys.FindRecordStartAid(dirFilename, int(total), createTime, filename, isDesc)
}

func LoadGeneralArticlesSameCreateTime(boardIDRaw *ptttype.BoardID_t, bid ptttype.Bid, startAid ptttype.Aid, endAid ptttype.Aid, createTime types.Time4) (summaries []*ptttype.ArticleSummaryRaw, startNumIdx ptttype.NumIdx, endNumIdx ptttype.NumIdx, err error) {

	total, err := cache.GetBTotalWithRetry(bid)
	if err != nil {
		return nil, 0, 0, err
	}
	if total == 0 {
		return nil, 0, 0, nil
	}

	if endAid == 0 {
		endAid = ptttype.Aid(total)
	}

	//3. get records
	filename, err := setBDir(boardIDRaw)
	if err != nil {
		return nil, 0, 0, err
	}

	nArticles := int(endAid - startAid + 1)

	maxAid := ptttype.Aid(total)
	summaries, err = cmsys.GetRecords(boardIDRaw, filename, startAid, nArticles, false, maxAid)
	if err != nil {
		return nil, 0, 0, err
	}
	if len(summaries) == 0 {
		return nil, 0, 0, nil
	}

	//filter with same create-time.
	newSummaries := make([]*ptttype.ArticleSummaryRaw, 0, len(summaries))
	startIdx := ptttype.NumIdx(startAid)
	numIdx := ptttype.NumIdx(0)
	for idx, each := range summaries {
		eachCreateTime, _ := each.Filename.CreateTime()
		if eachCreateTime != createTime {
			continue
		}

		numIdx = startIdx + ptttype.NumIdx(idx)

		if len(newSummaries) == 0 {
			startNumIdx = numIdx
		}

		newSummaries = append(newSummaries, each)
	}
	endNumIdx = numIdx

	return newSummaries, startNumIdx, endNumIdx, nil
}
