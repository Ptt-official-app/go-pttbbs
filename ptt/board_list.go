package ptt

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//LoadBoardSummary
//
//Load General Board Summary
//
//Params:
//	user
// 	uid
//	bid
//
//Return:
//	summary
//	err
func LoadBoardSummary(user *ptttype.UserecRaw, uid ptttype.Uid, bid ptttype.Bid) (summary *ptttype.BoardSummaryRaw, err error) {
	nBoardsInCache := cache.NumBoards()
	var boardStat *ptttype.BoardStat
	for currentIdx := ptttype.SortIdx(0); ; currentIdx++ {
		if int32(currentIdx) >= nBoardsInCache {
			break
		}

		eachBoardStat, err := loadGeneralBoardStat(user, uid, currentIdx, []byte(""))
		if err != nil {
			continue
		}
		if eachBoardStat == nil {
			continue
		}

		if eachBoardStat.Bid == bid {
			boardStat = eachBoardStat
			break
		}

	}

	if boardStat == nil {
		return nil, err
	}

	summary = parseBoardSummary(user, uid, boardStat)

	return summary, nil
}

//LoadGeneralBoards
//
//Load general boards by name.
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1142
//Params:
//	user
//	uid
//	startIdx: the idx in bsorted.
//  nBoards: try to get at most nBoards
//	keyword
//
//Return:
//  summary
//	nextIdx: next idx in bsorted.
//	err
func LoadGeneralBoards(user *ptttype.UserecRaw, uid ptttype.Uid, startIdx ptttype.SortIdx, nBoards int, keyword []byte) (summary []*ptttype.BoardSummaryRaw, nextIdx ptttype.SortIdx, err error) {

	nBoardsInCache := cache.NumBoards()

	boardStats := make([]*ptttype.BoardStat, 0, nBoards)
	currentIdx := ptttype.SortIdx(0)
	for currentIdx = startIdx; ; currentIdx++ {
		if int32(currentIdx) >= nBoardsInCache || len(boardStats) >= nBoards {
			break
		}

		eachBoardStat, err := loadGeneralBoardStat(user, uid, currentIdx, keyword)
		if err != nil {
			continue
		}
		if eachBoardStat == nil {
			continue
		}

		boardStats = append(boardStats, eachBoardStat)
	}

	summary, err = showBoardList(user, uid, boardStats)
	if err != nil {
		return nil, -1, err
	}

	if int32(currentIdx) == nBoardsInCache {
		currentIdx = -1
	}

	return summary, currentIdx, nil
}

//loadGeneralBoardStat
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1147
func loadGeneralBoardStat(user *ptttype.UserecRaw, uid ptttype.Uid, idx ptttype.SortIdx, keyword []byte) (*ptttype.BoardStat, error) {
	var bidInCache ptttype.BidInStore

	const bsort0sz = unsafe.Sizeof(cache.Shm.Raw.BSorted[0])
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BSorted)+bsort0sz*uintptr(ptttype.BSORT_BY_NAME)+uintptr(idx)*ptttype.BID_IN_STORE_SZ,
		ptttype.BID_IN_STORE_SZ,
		unsafe.Pointer(&bidInCache),
	)
	if bidInCache < 0 {
		return nil, nil
	}

	board := &ptttype.BoardHeaderRaw{}
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache),
		ptttype.BOARD_HEADER_RAW_SZ,
		unsafe.Pointer(board),
	)

	bid := bidInCache.ToBid()
	isGroupOp := groupOp(user, board)
	state := boardPermStat(user, uid, board, bid)
	if (board.Brdname[0] == '\x00') ||
		(board.BrdAttr&(ptttype.BRD_GROUPBOARD|ptttype.BRD_SYMBOLIC) != 0) ||
		!((state != ptttype.NBRD_INVALID) || isGroupOp) ||
		keywordNotInTitle(&board.Title, keyword) {
		return nil, nil
	}

	boardStat := newBoardStat(bidInCache, state, board, isGroupOp)

	return boardStat, nil
}

//newBoardStat
func newBoardStat(bidInCache ptttype.BidInStore, state ptttype.BoardStatAttr, board *ptttype.BoardHeaderRaw, isGroupOp bool) (boardStat *ptttype.BoardStat) {
	boardStat = &ptttype.BoardStat{}

	boardStat.Bid = bidInCache.ToBid()
	boardStat.Attr = state

	boardStat.Board = board
	boardStat.IsGroupOp = isGroupOp

	//XXX need to modify this by having state with NBRD_SET_POSTMASK
	//XXX this is a hack to ensure the brd-postmask
	var brd_postmask = ptttype.BRD_POSTMASK
	if (board.BrdAttr&ptttype.BRD_HIDE != 0) && (board.BrdAttr&ptttype.BRD_POSTMASK == 0) && state == ptttype.NBRD_BOARD {
		cache.Shm.SetOrUint32(
			unsafe.Offsetof(cache.Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache)+ptttype.BOARD_HEADER_BRD_ATTR_OFFSET,
			uint32(brd_postmask),
		)
		board.BrdAttr |= brd_postmask
	}

	return boardStat
}

//keywordNotInTitle
//
//TITLE_MATCH in board.c
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L14
func keywordNotInTitle(title *ptttype.BoardTitle_t, keyword []byte) bool {
	result := (len(keyword) > 0) && (types.Cstrcasestr(title[:], keyword) < 0)

	return result
}

//showBoardList
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1409
func showBoardList(user *ptttype.UserecRaw, uid ptttype.Uid, boardStats []*ptttype.BoardStat) (summary []*ptttype.BoardSummaryRaw, err error) {
	summary = make([]*ptttype.BoardSummaryRaw, len(boardStats))
	for idx, eachStat := range boardStats {
		summary[idx] = parseBoardSummary(user, uid, eachStat)
	}

	return summary, nil
}

//parseBoardSummary
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1460
func parseBoardSummary(user *ptttype.UserecRaw, uid ptttype.Uid, boardStat *ptttype.BoardStat) (summary *ptttype.BoardSummaryRaw) {

	//XXX we do not deal with fav in go-bbs.
	if boardStat.Attr&ptttype.NBRD_LINE != 0 {
		return &ptttype.BoardSummaryRaw{Bid: boardStat.Bid, StatAttr: boardStat.Attr}
	}

	//XXX we do not deal with fav in go-bbs.
	if boardStat.Attr&ptttype.NBRD_FOLDER != 0 {
		return &ptttype.BoardSummaryRaw{Bid: boardStat.Bid, StatAttr: boardStat.Attr}
	}

	//hidden board
	board := boardStat.Board
	if !boardStat.IsGroupOp && boardStat.Attr == ptttype.NBRD_INVALID {
		reason := ptttype.RESTRICT_REASON_FORBIDDEN
		if board.BrdAttr&ptttype.BRD_HIDE != 0 {
			reason = ptttype.RESTRICT_REASON_HIDDEN
		}
		summary = &ptttype.BoardSummaryRaw{
			Bid:      boardStat.Bid,
			BrdAttr:  board.BrdAttr,
			StatAttr: boardStat.Attr,
			Brdname:  &board.Brdname,
			Reason:   reason,
		}
		if ptttype.USE_REAL_DESC_FOR_HIDDEN_BOARD_IN_MYFAV {
			summary.Title = &board.Title
		}

		return summary
	}

	bidInCache := boardStat.Bid.ToBidInStore()
	var lastPostTime types.Time4
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.LastPostTime)+types.TIME4_SZ*uintptr(bidInCache),
		types.TIME4_SZ,
		unsafe.Pointer(&lastPostTime),
	)

	var total int32
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.Total)+types.INT32_SZ*uintptr(bidInCache),
		types.INT32_SZ,
		unsafe.Pointer(&total),
	)

	summary = &ptttype.BoardSummaryRaw{
		Bid:          boardStat.Bid,
		BrdAttr:      board.BrdAttr,
		StatAttr:     boardStat.Attr,
		Brdname:      &board.Brdname,
		Title:        &board.Title,
		BM:           board.BM.ToBMs(),
		LastPostTime: lastPostTime,
		NUser:        board.NUser,
		Total:        total,
	}

	return summary
}

//groupOp
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1579
func groupOp(user *ptttype.UserecRaw, board *ptttype.BoardHeaderRaw) bool {
	if user.UserLevel.HasUserPerm(ptttype.PERM_NOCITIZEN) {
		return false
	}

	if user.UserLevel.HasUserPerm(ptttype.PERM_BOARD) {
		return true
	}

	if is_uBM(&user.UserID, &board.BM) {
		return true
	}

	return false
}
