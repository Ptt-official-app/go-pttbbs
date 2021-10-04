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
		startIdx, err := ptt.FindArticleStartIdx(userecRaw, uid, boardIDRaw, bid, createTime, filename, false)
		result = append(result, startIdx)
		// TODO is need recover deleted items if get error?
		err = ptt.DeleteArticles(boardIDRaw, filename)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
