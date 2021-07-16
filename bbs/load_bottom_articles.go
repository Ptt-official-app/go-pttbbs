package bbs

import "github.com/Ptt-official-app/go-pttbbs/ptt"

func LoadBottomArticles(uuserID UUserID, bboardID BBoardID) (summaries []*ArticleSummary, err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, err
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	summariesRaw, err := ptt.LoadBottomArticles(userecRaw, uid, boardIDRaw, bid)
	if err != nil {
		return nil, err
	}

	summaries = make([]*ArticleSummary, len(summariesRaw))
	for idx, each := range summariesRaw {
		eachSummary := NewArticleSummaryFromRaw(bboardID, each)
		summaries[idx] = eachSummary
	}

	return summaries, nil
}
