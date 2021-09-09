package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type BoardDetail struct {
	Brdname            string          `json:"brdname"`
	RealTitle          []byte          `json:"title"` //Require to separate RealTitle, BoardClass, BoardType, because it's hard to parse in utf8
	BM                 []UUserID       `json:"moderators"`
	BrdAttr            ptttype.BrdAttr `json:"attr"`
	Gid                ptttype.Bid     `json:"pttgid"`
	Bid                ptttype.Bid     `json:"pttbid"`
	BBoardID           BBoardID        `json:"bid"`
	ChessCountry       byte            `json:"chesscountry"`
	VoteLimitLogins    uint8           `json:"votelimitlogins"`
	BUpdate            types.Time4     `json:"bupdate"`
	PostLimitLogins    uint8           `json:"postlimitlogins"`
	BVote              uint8           `json:"bvote"`
	VTime              types.Time4     `json:"vtime"`
	Level              ptttype.PERM    `json:"level"`
	PermReload         types.Time4     `json:"permreload"`
	Next               [2]ptttype.Bid  `json:"next"`
	FirstChild         [2]ptttype.Bid  `json:"firstchild"`
	Parent             ptttype.Bid     `json:"parent"`
	ChildCount         int32           `json:"childcount"`
	NUser              int32           `json:"nuser"`
	PostExpire         int32           `json:"postexpire"`
	EndGamble          types.Time4     `json:"endgamble"`
	PostType           []byte          `json:"posttype"`
	FastRecommendPause uint8           `json:"fastrecommendpause"`
	VoteLimitBadpost   uint8           `json:"votelimitbadpost"`
}

func NewBoardDetailFromRaw(boardHeaderRaw *ptttype.BoardHeaderRaw, bid ptttype.Bid) *BoardDetail {
	bmsRaw := boardHeaderRaw.BM.ToBMs()
	bms := make([]UUserID, len(bmsRaw))
	for idx, each := range bmsRaw {
		bms[idx] = UUserID(types.CstrToString(each[:]))
	}

	boardDetail := &BoardDetail{
		Brdname:            types.CstrToString(boardHeaderRaw.Brdname[:]),
		RealTitle:          types.CstrToBytes(boardHeaderRaw.Title[7:]),
		BM:                 bms,
		BrdAttr:            boardHeaderRaw.BrdAttr,
		Gid:                boardHeaderRaw.Gid,
		Bid:                bid,
		BBoardID:           ToBBoardID(bid, &boardHeaderRaw.Brdname),
		ChessCountry:       byte(boardHeaderRaw.ChessCountry),
		VoteLimitLogins:    boardHeaderRaw.VoteLimitLogins,
		BUpdate:            boardHeaderRaw.BUpdate,
		PostLimitLogins:    boardHeaderRaw.VoteLimitLogins,
		BVote:              boardHeaderRaw.BVote,
		VTime:              boardHeaderRaw.VTime,
		Level:              boardHeaderRaw.Level,
		PermReload:         boardHeaderRaw.PermReload,
		Next:               [2]ptttype.Bid{},
		FirstChild:         [2]ptttype.Bid{},
		Parent:             boardHeaderRaw.Parent,
		ChildCount:         boardHeaderRaw.ChildCount,
		NUser:              boardHeaderRaw.ChildCount,
		PostExpire:         boardHeaderRaw.PostExpire,
		EndGamble:          boardHeaderRaw.EndGamble,
		PostType:           types.CstrToBytes(boardHeaderRaw.PostType[:]),
		FastRecommendPause: boardHeaderRaw.FastRecommendPause,
		VoteLimitBadpost:   boardHeaderRaw.VoteLimitBadpost,
	}

	return boardDetail
}
