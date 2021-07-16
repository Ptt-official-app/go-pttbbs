package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func GetArticle(uuserID UUserID, bboardID BBoardID, articleID ArticleID, retrieveTS types.Time4) (content []byte, mtime types.Time4, err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, 0, ErrInvalidParams
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, 0, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, 0, ErrInvalidParams
	}

	filename := articleID.ToFilename()

	return ptt.ReadPost(userecRaw, uid, boardIDRaw, bid, filename, retrieveTS)
}
