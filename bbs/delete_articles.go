package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func DeleteArticles(uuserID UUserID, bboardID BBoardID, articleIDs []ArticleID, ip string) ([]ptttype.SortIdx, error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, err
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, err
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, err
	}

	var result []ptttype.SortIdx
	for _, articleID := range articleIDs {
		filename := articleID.ToFilename()
		createTime, err := filename.CreateTime()
		if err != nil {
			return nil, err
		}
		startIdx, err := ptt.FindArticleStartIdx(userecRaw, uid, boardIDRaw, bid, createTime, filename, false)
		if err != nil {
			return nil, err
		}
		// FindArticleStartIdx only find nearest idx, so we must make sure filename is exactly correct
		summariesRaw, _, _, _, _ := ptt.LoadGeneralArticles(userecRaw, uid, boardIDRaw, bid, startIdx, 1, true)
		if len(summariesRaw) == 1 {
			articleSummary := NewArticleSummaryFromRaw(bboardID, summariesRaw[0])
			if articleID == articleSummary.ArticleID {
				err = ptt.DeleteArticles(boardIDRaw, filename, startIdx)
				// TODO is need recover deleted items if get error?
				if err != nil {
					return nil, err
				}
				result = append(result, startIdx)
			}
		}
	}
	return result, nil
}
