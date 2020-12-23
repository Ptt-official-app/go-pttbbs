package ptttype

import "github.com/Ptt-official-app/go-pttbbs/types"

type BoardSummaryRaw struct {
	Bid          Bid
	BrdAttr      BrdAttr
	StatAttr     BoardStatAttr
	Brdname      *BoardID_t
	Title        *BoardTitle_t
	BM           []*UserID_t
	Reason       RestrictReason
	LastPostTime types.Time4
	NUser        int32
	Total        int32
}
