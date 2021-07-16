package ptttype

import "github.com/Ptt-official-app/go-pttbbs/types"

type BoardSummaryRaw struct {
	Gid      Bid
	Bid      Bid
	BrdAttr  BrdAttr
	StatAttr BoardStatAttr
	Brdname  *BoardID_t

	Title        *BoardTitle_t
	BM           []*UserID_t
	Reason       RestrictReason
	LastPostTime types.Time4
	NUser        int32
	Total        int32
}

func NewBoardSummaryRawWithReason(boardStat *BoardStat) (summary *BoardSummaryRaw) {
	board := boardStat.Board

	reason := RESTRICT_REASON_FORBIDDEN
	if board.BrdAttr&BRD_HIDE != 0 {
		reason = RESTRICT_REASON_HIDDEN
	}

	summary = &BoardSummaryRaw{
		Gid:      board.Gid,
		Bid:      boardStat.Bid,
		BrdAttr:  board.BrdAttr,
		StatAttr: boardStat.Attr,
		Brdname:  &board.Brdname,

		Reason: reason,
	}
	if USE_REAL_DESC_FOR_HIDDEN_BOARD_IN_MYFAV {
		summary.Title = &board.Title
	}

	return summary
}

func NewBoardSummaryRaw(boardStat *BoardStat, lastPostTime types.Time4, total int32) (summary *BoardSummaryRaw) {
	board := boardStat.Board
	summary = &BoardSummaryRaw{
		Gid:      board.Gid,
		Bid:      boardStat.Bid,
		BrdAttr:  board.BrdAttr,
		StatAttr: boardStat.Attr,
		Brdname:  &board.Brdname,

		Title:        &board.Title,
		BM:           board.BM.ToBMs(),
		LastPostTime: lastPostTime,
		NUser:        board.NUser,
		Total:        total,
	}

	return summary
}
