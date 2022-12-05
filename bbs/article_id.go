package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type ArticleID string

// ToArticleID
//
// aidc is with fixed-size (8 bytes), no need the separator to do the separation.
func ToArticleID(filename *ptttype.Filename_t) ArticleID {
	aidc := filename.ToAidu().ToAidc()
	aidcStr := types.CstrToString(aidc[:])
	return ArticleID(aidcStr)
}

// ToRaw
//
// 1st 8 bytes are aidc.
// bytes starting 8th are ownerID
func (a ArticleID) ToRaw() (filename *ptttype.Filename_t) {
	// 1st 8 bytes are aidc
	aidc := &ptttype.Aidc{}
	copy(aidc[:], []byte(a[:8]))
	filename = aidc.ToAidu().ToFN()

	return filename
}

func (a *ArticleID) ToFilename() (filename *ptttype.Filename_t) {
	filename = a.ToRaw()
	return filename
}
