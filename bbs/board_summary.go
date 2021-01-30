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

	bms := make([]UUserID, len(boardSummaryRaw.BM))
	for idx, each := range boardSummaryRaw.BM {
		bms[idx] = UUserID(types.CstrToString(each[:]))
	}
	boardSummary := &BoardSummary{
		BBoardID:     ToBBoardID(boardSummaryRaw.Bid, boardSummaryRaw.Brdname),
		BrdAttr:      boardSummaryRaw.BrdAttr,
		StatAttr:     boardSummaryRaw.StatAttr,
		Brdname:      types.CstrToString(boardSummaryRaw.Brdname[:]),
		BoardClass:   types.CstrToBytes(boardSummaryRaw.Title[:4]),
		BoardType:    types.CstrToBytes(boardSummaryRaw.Title[5:7]),
		RealTitle:    types.CstrToBytes(boardSummaryRaw.Title[7:]),
		BM:           bms,
		Reason:       boardSummaryRaw.Reason,
		LastPostTime: boardSummaryRaw.LastPostTime,
		Total:        boardSummaryRaw.Total,
		NUser:        boardSummaryRaw.NUser,
	}

	return boardSummary
}
