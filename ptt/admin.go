package ptt

import (
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func addBoardRecord(board *ptttype.BoardHeaderRaw) (bid ptttype.Bid, err error) {
	emptyBrdname := &ptttype.BoardID_t{}
	bid, err = cache.GetBid(emptyBrdname)
	if err != nil {
		return 0, err
	}

	// able to add the board to existing empty board.
	if bid.IsValid() {
		err = cmsys.SubstituteRecord(ptttype.FN_BOARD, board, ptttype.BOARD_HEADER_RAW_SZ, int32(bid))
		if err != nil {
			return 0, err
		}
		cache.ResetBoard(bid)
		cache.SortBCache()

		return bid, nil
	}

	// add board to the end of the board-record
	nBoards := cache.NumBoards()
	if nBoards >= ptttype.MAX_BOARD {
		return 0, ErrTooManyBoards
	}

	_, err = cmsys.AppendRecord(ptttype.FN_BOARD, board, ptttype.BOARD_HEADER_RAW_SZ)
	if err != nil {
		return 0, err
	}

	bid, err = cache.AddbrdTouchCache()

	return bid, err
}

//mNewbrd
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/admin.c#L1071
func mNewbrd(
	user *ptttype.UserecRaw,
	clsBid ptttype.Bid,
	brdname *ptttype.BoardID_t,
	brdClass []byte,
	brdTitle []byte,
	bms *ptttype.BM_t,
	brdAttr ptttype.BrdAttr,
	level ptttype.PERM,
	chessCountry ptttype.ChessCode,
	isGroup bool,
	isRecover bool,
) (board *ptttype.BoardHeaderRaw, bid ptttype.Bid, err error) {

	if !clsBid.IsValid() {
		return nil, 0, ptttype.ErrInvalidBid
	}

	if !brdname.IsValid() {
		return nil, 0, ptttype.ErrInvalidBoardID
	}

	bid, err = cache.GetBid(brdname)
	if err == nil && bid > 0 {
		return nil, 0, ptttype.ErrBoardIDAlreadyExists
	}

	// mkdir
	dirname := path.SetBPath(brdname)
	err = types.Mkdir(dirname)
	if os.IsExist(err) && isRecover {
		err = nil
	}
	if err != nil {
		return nil, 0, err
	}

	// BMs
	parsedBMs := cache.SanitizeBMs(bms)

	// title
	title := &ptttype.BoardTitle_t{}
	copy(title[:4], brdClass)
	copy(title[7:], brdTitle)
	title[4] = ' '
	if isGroup {
		copy(title[5:7], ptttype.BRD_SYMBOL_GROUP)
	} else {
		copy(title[5:7], ptttype.BRD_SYMBOL_BOARD)
	}

	// brdAttr
	if ptttype.DEFAULT_AUTOCPLOG {
		brdAttr |= ptttype.BRD_CPLOG
	}

	if isGroup {
		brdAttr |= ptttype.BRD_GROUPBOARD
		brdAttr &= ^ptttype.BRD_CPLOG
	} else {
		brdAttr &= ^ptttype.BRD_GROUPBOARD
	}

	// level
	if !user.UserLevel.HasUserPerm(ptttype.PERM_BOARD) || brdAttr&ptttype.BRD_HIDE != 0 {
		brdAttr &= ^ptttype.BRD_POSTMASK
		level = 0
	}

	// board
	board = &ptttype.BoardHeaderRaw{
		Brdname:      *brdname,
		Title:        *title,
		BM:           *parsedBMs,
		BrdAttr:      brdAttr,
		Level:        level,
		ChessCountry: chessCountry,
		Gid:          clsBid,
	}

	bid, err = addBoardRecord(board)
	if err != nil {
		return nil, 0, err
	}

	board, err = cache.GetBCache(bid)
	if err != nil {
		return nil, bid, err
	}

	return board, bid, nil
}
