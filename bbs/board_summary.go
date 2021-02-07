package bbs

import (
	"encoding/base64"
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type BoardSummary struct {
	Gid          ptttype.Bid            `json:"pttgid"`
	Bid          ptttype.Bid            `json:"pttbid"`
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

	IdxByName  string `json:"idx_name"`
	IdxByClass string `json:"idx_class"`
}

func NewBoardSummaryFromRaw(boardSummaryRaw *ptttype.BoardSummaryRaw) *BoardSummary {

	bms := make([]UUserID, len(boardSummaryRaw.BM))
	for idx, each := range boardSummaryRaw.BM {
		bms[idx] = UUserID(types.CstrToString(each[:]))
	}
	boardSummary := &BoardSummary{
		Gid:          boardSummaryRaw.Gid,
		Bid:          boardSummaryRaw.Bid,
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
	boardSummary.IdxByName = serializeBoardIdxByNameStr(boardSummary)
	boardSummary.IdxByClass = serializeBoardIdxByClassStr(boardSummary)

	return boardSummary
}

func serializeBoardIdxByNameStr(summary *BoardSummary) (idxStr string) {
	return string(summary.Brdname)
}

func deserializeBoardIdxByNameStr(idxStr string) (brdname string, err error) {
	return idxStr, nil
}

func serializeBoardIdxByClassStr(summary *BoardSummary) (idxStr string) {
	return base64.RawURLEncoding.EncodeToString(summary.BoardClass) + "@" + summary.Brdname
}

func deserializeBoardIdxByClassStr(idxStr string) (boardClass []byte, brdname string, err error) {
	if idxStr == "" {
		return nil, "", nil
	}
	theList := strings.Split(idxStr, "@")
	if len(theList) != 2 {
		return nil, "", ErrInvalidParams
	}

	boardClass, err = base64.RawURLEncoding.DecodeString(theList[0])
	if err != nil {
		return nil, "", err
	}

	return boardClass, theList[1], nil
}
