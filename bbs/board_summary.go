package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type BoardSummary struct {
	BBoardID     BBoardID              `json:"bid"`
	BrdAttr      ptttype.BrdAttr       `json:"brdAttr"`
	StatAttr     ptttype.BoardStatAttr `json:"statAttr"`
	Brdname      string                `json:"brdname"`
	RealTitle    string                `json:"title"` //Require to separate RealTitle, BoardClass, BoardType, because it's hard to parse in utf8
	BoardClass   string                `json:"class"`
	BoardType    string                `json:"boardType"` //□, ◎, Σ
	BM           []string              `json:"moderators"`
	Reason       string                `json:"reason"`
	LastPostTime types.Time4           `json:"lastPostTime"`
	NUser        int32                 `json:"nUser"`
	Total        int32                 `json:"total"`
}

func NewBoardSummaryFromRaw(boardSummaryRaw *ptttype.BoardSummaryRaw) *BoardSummary {

	boardSummary := &BoardSummary{}

	boardSummary.BBoardID = ToBBoardID(boardSummaryRaw.Bid, boardSummaryRaw.Brdname)
	boardSummary.BrdAttr = boardSummaryRaw.BrdAttr
	boardSummary.StatAttr = boardSummaryRaw.StatAttr
	boardSummary.Brdname = types.CstrToString(boardSummaryRaw.Brdname[:])

	boardSummary.BoardClass = types.Big5ToUtf8(boardSummaryRaw.Title[:4])
	boardSummary.BoardType = types.Big5ToUtf8(boardSummaryRaw.Title[5:7])
	boardSummary.RealTitle = types.Big5ToUtf8(boardSummaryRaw.Title[7:])
	boardSummary.BM = make([]string, len(boardSummaryRaw.BM))
	if len(boardSummaryRaw.BM) > 0 {
		for idx, each := range boardSummaryRaw.BM {
			boardSummary.BM[idx] = types.CstrToString(each[:])
		}
	}
	boardSummary.Reason = boardSummaryRaw.Reason.String()
	boardSummary.LastPostTime = boardSummaryRaw.LastPostTime
	boardSummary.Total = boardSummaryRaw.Total

	return boardSummary
}
