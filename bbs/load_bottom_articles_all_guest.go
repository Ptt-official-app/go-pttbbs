package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func LoadBottomArticlesAllGuest(bboardID BBoardID) (summaries []*ArticleSummary, err error) {
	boardIDRaw, err := bboardID.ToRawAllGuest()
	if err != nil {
		return nil, err
	}

	summariesRaw, err := ptt.LoadBottomArticlesAllGuest(boardIDRaw)
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
