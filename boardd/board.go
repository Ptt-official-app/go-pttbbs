package boardd

import "github.com/Ptt-official-app/go-pttbbs/ptttype"

func boardToBoardSummaryRaw(board *Board) (summary *ptttype.BoardSummaryRaw) {
	brdname := &ptttype.BoardID_t{}
	copy(brdname[:], []byte(board.Name))

	title := &ptttype.BoardTitle_t{}
	copy(title[:], []byte(board.Title))
	bms := &ptttype.BM_t{}
	copy(bms[:], []byte(board.RawModerators))

	return &ptttype.BoardSummaryRaw{
		Gid:     ptttype.Bid(board.Parent),
		Bid:     ptttype.Bid(board.Bid),
		BrdAttr: ptttype.BrdAttr(board.Attributes),
		Brdname: brdname,

		Title: title,
		BM:    bms.ToBMs(),
		NUser: int32(board.NumUsers),
		Total: int32(board.NumPosts),
	}
}
