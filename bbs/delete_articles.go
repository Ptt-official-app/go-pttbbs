package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/sirupsen/logrus"
)

func DeleteArticles(uuserID UUserID, bboardID BBoardID, articleIDs []ArticleID, ip string) error {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return err
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return err
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return err
	}
	for _, articleID := range articleIDs {
		filename := articleID.ToRaw()
		startIdx, err := ptt.FindArticleStartIdx(userecRaw, uid, boardIDRaw, bid, 0, filename, false)
		logrus.Info(startIdx)
		if err != nil {
			return err
		}
	}
	return nil
}
