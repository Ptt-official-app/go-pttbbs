package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
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
func LoadGeneralArticles(user *ptttype.UserecRaw, uid ptttype.Uid, boardIDRaw *ptttype.BoardID_t, bid ptttype.Bid, startIdx ptttype.SortIdx, nArticles int) (summaryRaw []*ptttype.ArticleSummaryRaw, nextIdx ptttype.SortIdx, isNewest bool, err error) {

	//1. check perm.
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, -1, false, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, -1, false, ErrNotPermitted
	}

	//2. bcache preparation.
	total := cache.GetBTotal(bid)
	if total == 0 {
		err = cache.SetBTotal(bid)
		if err != nil {
			return nil, -1, false, err
		}
		err = cache.SetBottomTotal(bid)
		if err != nil {
			return nil, -1, false, err
		}

		total = cache.GetBTotal(bid)
		if total == 0 { //no data
			return nil, -1, false, nil
		}
	}

	//3. get records
	filename, err := setBDir(boardIDRaw)
	if err != nil {
		return nil, -1, false, err
	}
	//3.1. ensure recordStartAid and startAid
	if startIdx == 0 {
		startIdx = ptttype.SortIdx(total)
	}
	startAid := ptttype.Aid(startIdx)
	if startAid > ptttype.Aid(total) {
		startAid = ptttype.Aid(total)
	}

	recordStartAid := startAid - ptttype.Aid(nArticles) + 1 //startAid is included.
	if recordStartAid < 1 {
		recordStartAid = 1
	}
	nArticles = int(startAid - recordStartAid + 1) //startAid is included

	if nArticles < 1 {
		return nil, -1, false, ErrInvalidParams
	}

	summaryRaw, err = cmsys.GetRecords(boardIDRaw, filename, recordStartAid, nArticles)
	if err != nil {
		return nil, -1, false, err
	}

	//4. return
	nextIdx = -1
	if recordStartAid > 1 {
		nextIdx = ptttype.SortIdx(recordStartAid) - 1
	}

	isNewest = startAid == ptttype.Aid(total)
	return summaryRaw, nextIdx, isNewest, nil
}
