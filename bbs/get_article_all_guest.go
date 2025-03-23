package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func GetArticleAllGuest(bboardID BBoardID, articleID ArticleID, retrieveTS types.Time4, isHash bool) (content []byte, mtime types.Time4, sum cmsys.Fnv64_t, err error) {
	boardIDRaw, err := bboardID.ToRawAllGuest()
	if err != nil {
		return nil, 0, 0, err
	}

	filename := articleID.ToFilename()

	return ptt.ReadPostAllGuest(boardIDRaw, filename, retrieveTS, isHash)
}
