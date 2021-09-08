package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
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
func LoadGeneralArticles(user *ptttype.UserecRaw, uid ptttype.UID, boardIDRaw *ptttype.BoardID_t, bid ptttype.Bid, startIdx ptttype.SortIdx, nArticles int, isDesc bool) (summaries []*ptttype.ArticleSummaryRaw, isNewest bool, nextSummary *ptttype.ArticleSummaryRaw, startNumIdx ptttype.SortIdx, err error) {
	// 1. check perm.
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, false, nil, 0, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, false, nil, 0, ErrNotPermitted
	}

	// 2. bcache preparation.
	total, err := cache.GetBTotalWithRetry(bid)
	if err != nil {
		return nil, false, nil, 0, err
	}
	if total == 0 {
		return nil, true, nil, 0, nil
	}

	// 3. get records
	filename, err := setBDir(boardIDRaw)
	if err != nil {
		return nil, false, nil, 0, err
	}
	// 3.1. ensure startAid
	if startIdx == 0 && isDesc {
		startIdx = ptttype.SortIdx(total)
	}

	summaries, err = cmsys.GetRecords(boardIDRaw, filename, startIdx, nArticles+1, isDesc)
	if err != nil {
		return nil, false, nil, 0, err
	}

	// 4. return
	isNewest = false
	if isDesc {
		isNewest = startIdx == ptttype.SortIdx(total)
	} else {
		isNewest = len(summaries) != nArticles+1
	}

	if len(summaries) == nArticles+1 {
		nextSummary = summaries[nArticles]
		summaries = summaries[:nArticles]
	}

	return summaries, isNewest, nextSummary, startIdx, nil
}

func LoadBottomArticles(user *ptttype.UserecRaw, uid ptttype.UID, boardIDRaw *ptttype.BoardID_t, bid ptttype.Bid) (summaries []*ptttype.ArticleSummaryRaw, err error) {
	// 1. check perm.
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, ErrNotPermitted
	}

	// 2. bcache preparation.
	total := cache.GetBottomTotal(bid)
	if total == 0 {
		return nil, nil
	}

	// 3. get records
	filename, err := path.SetBFile(boardIDRaw, ptttype.FN_DIR_BOTTOM)
	if err != nil {
		return nil, err
	}

	summaries, err = cmsys.GetRecords(boardIDRaw, filename, 1, int(total), false)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func FindArticleStartIdx(user *ptttype.UserecRaw, uid ptttype.UID, boardID *ptttype.BoardID_t, bid ptttype.Bid, createTime types.Time4, filename *ptttype.Filename_t, isDesc bool) (startIdx ptttype.SortIdx, err error) {
	// 1. check perm.
	board, err := cache.GetBCache(bid)
	if err != nil {
		return -1, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return -1, ErrNotPermitted
	}

	// 3. get records
	dirFilename, err := setBDir(boardID)
	if err != nil {
		return -1, err
	}

	total, err := cache.GetBTotalWithRetry(bid)
	if err != nil {
		return -1, err
	}
	if total == 0 {
		return -1, ErrNoRecord
	}

	return cmsys.FindRecordStartIdx(dirFilename, int(total), createTime, filename, isDesc)
}

func LoadGeneralArticlesSameCreateTime(boardIDRaw *ptttype.BoardID_t, bid ptttype.Bid, startIdx ptttype.SortIdx, endIdx ptttype.SortIdx, createTime types.Time4) (summaries []*ptttype.ArticleSummaryRaw, startNumIdx ptttype.SortIdx, endNumIdx ptttype.SortIdx, err error) {
	total, err := cache.GetBTotalWithRetry(bid)
	if err != nil {
		return nil, 0, 0, err
	}
	if total == 0 {
		return nil, 0, 0, nil
	}

	if endIdx == 0 {
		endIdx = ptttype.SortIdx(total)
	}

	// 3. get records
	filename, err := setBDir(boardIDRaw)
	if err != nil {
		return nil, 0, 0, err
	}

	nArticles := int(endIdx - startIdx + 1)

	summaries, err = cmsys.GetRecords(boardIDRaw, filename, startIdx, nArticles, false)
	if err != nil {
		return nil, 0, 0, err
	}
	if len(summaries) == 0 {
		return nil, 0, 0, nil
	}

	// filter with same create-time.
	newSummaries := make([]*ptttype.ArticleSummaryRaw, 0, len(summaries))
	startNumIdx = ptttype.SortIdx(startIdx)
	numIdx := ptttype.SortIdx(0)
	for idx, each := range summaries {
		eachCreateTime, _ := each.Filename.CreateTime()
		if eachCreateTime != createTime {
			continue
		}

		numIdx = startNumIdx + ptttype.SortIdx(idx)

		if len(newSummaries) == 0 {
			startNumIdx = numIdx
		}

		newSummaries = append(newSummaries, each)
	}
	if len(newSummaries) == 0 {
		return nil, 0, 0, nil
	}

	endNumIdx = numIdx

	return newSummaries, startNumIdx, endNumIdx, nil
}
