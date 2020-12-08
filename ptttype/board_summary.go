package ptttype

import "github.com/Ptt-official-app/go-pttbbs/types"

type RestrictReason uint8

const (
	_ RestrictReason = iota
	RESTRICT_REASON_FORBIDDEN
	RESTRICT_REASON_HIDDEN
)

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
