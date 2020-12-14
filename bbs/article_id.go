package bbs

import (
	"strconv"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type ArticleID string

func ToArticleID(aid ptttype.Aid, filename ptttype.Filename_t) ArticleID {
	aidc := filename.ToAidu().ToAidc()
	aidcStr := types.CstrToString(aidc[:])
	return ArticleID(strconv.Itoa(int(aid)) + "_" + aidcStr)
}
