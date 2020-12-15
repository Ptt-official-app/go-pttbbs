package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type ArticleID string

func ToArticleID(filename ptttype.Filename_t) ArticleID {
	aidc := filename.ToAidu().ToAidc()
	aidcStr := types.CstrToString(aidc[:])
	return ArticleID(aidcStr)
}
