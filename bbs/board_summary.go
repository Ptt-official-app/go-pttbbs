package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type BoardSummary struct {
	BBoardID     BBoardID               `json:"bid"`
	BrdAttr      ptttype.BrdAttr        `json:"attr"`
	StatAttr     ptttype.BoardStatAttr  `json:"user_attr"`
	Brdname      string                 `json:"brdname"`
	RealTitle    []byte                 `json:"title"` //Require to separate RealTitle, BoardClass, BoardType, because it's hard to parse in utf8
	BoardClass   []byte                 `json:"class"`
	BoardType    []byte                 `json:"type"` //□, ◎, Σ
	BM           []UUserID              `json:"moderators"`
	Reason       ptttype.RestrictReason `json:"reason"`
	LastPostTime types.Time4            `json:"last_post_time"`
	NUser        int32                  `json:"number_of_user"`
	Total        int32                  `json:"total"`
	Read         bool                   `json:"read"`
}

func NewBoardSummaryFromRaw(boardSummaryRaw *ptttype.BoardSummaryRaw) *BoardSummary {

	boardSummary := &BoardSummary{}

	boardSummary.BBoardID = ToBBoardID(boardSummaryRaw.Bid, boardSummaryRaw.Brdname)
	boardSummary.BrdAttr = boardSummaryRaw.BrdAttr
	boardSummary.StatAttr = boardSummaryRaw.StatAttr
	boardSummary.Brdname = types.CstrToString(boardSummaryRaw.Brdname[:])

	boardSummary.BoardClass = types.CstrToBytes(boardSummaryRaw.Title[:4])
	boardSummary.BoardType = types.CstrToBytes(boardSummaryRaw.Title[5:7])
	boardSummary.RealTitle = types.CstrToBytes(boardSummaryRaw.Title[7:])
	boardSummary.BM = make([]UUserID, len(boardSummaryRaw.BM))
	if len(boardSummaryRaw.BM) > 0 {
		for idx, each := range boardSummaryRaw.BM {
			boardSummary.BM[idx] = UUserID(types.CstrToString(each[:]))
		}
	}
	boardSummary.Reason = boardSummaryRaw.Reason
	boardSummary.LastPostTime = boardSummaryRaw.LastPostTime
	boardSummary.Total = boardSummaryRaw.Total

	return boardSummary
}
