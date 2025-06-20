package bbs

import "github.com/Ptt-official-app/go-pttbbs/cache"

func LoadBoardSummaryAllGuest(bboardID BBoardID) (summary *BoardSummary, err error) {
	boardIDRaw, err := bboardID.ToRawAllGuest()
	if err != nil {
		return nil, err
	}

	total := cache.GetBTotalAllGuest(boardIDRaw)

	lastPosttime, err := cache.GetLastPosttimeAllGuest(boardIDRaw)
	if err != nil {
		return nil, err
	}

	summary = &BoardSummary{
		BBoardID:     bboardID,
		Brdname:      bboardID.ToBrdnameAllGuest(),
		Total:        total,
		LastPostTime: lastPosttime,
	}

	return summary, nil
}
