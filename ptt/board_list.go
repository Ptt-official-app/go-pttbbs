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

	bidInCache := bid.ToBidInStore()

	if bidInCache < 0 {
		return nil, nil
	}
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, err
	}
	isGroupOp := groupOp(user, board)
	state := boardPermStat(user, uid, board, bid)
	boardStat := newBoardStat(bidInCache, state, board, isGroupOp)

	if boardStat == nil {
		return nil, err
	}

	summary = parseBoardSummary(user, uid, boardStat)

	return summary, nil
}

//LoadHotBoards
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1125
func LoadHotBoards(user *ptttype.UserecRaw, uid ptttype.Uid) (summary []*ptttype.BoardSummaryRaw, err error) {
	nBoards := cache.NHots()

	boardStats := make([]*ptttype.BoardStat, 0, nBoards)

	for idx := uint8(0); idx < nBoards; idx++ {
		eachBoardStat := loadHotBoardStat(user, uid, idx)

		if eachBoardStat == nil {
			continue
		}

		boardStats = append(boardStats, eachBoardStat)
	}

	summary, err = showBoardList(user, uid, boardStats)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

//loadHotBoardStat
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1147
func loadHotBoardStat(user *ptttype.UserecRaw, uid ptttype.Uid, idx uint8) *ptttype.BoardStat {

	//read bid-in-cache
	var bidInCache ptttype.BidInStore

	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.HBcache)+uintptr(idx)*ptttype.BID_IN_STORE_SZ,
		ptttype.BID_IN_STORE_SZ,
		unsafe.Pointer(&bidInCache),
	)
	if bidInCache < 0 {
		return nil
	}

	//get board
	board := &ptttype.BoardHeaderRaw{}
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache),
		ptttype.BOARD_HEADER_RAW_SZ,
		unsafe.Pointer(board),
	)

	//board-stat
	//assuming that the hot-boards can be accessed by the public.
	bid := bidInCache.ToBid()
	isGroupOp := groupOp(user, board)
	state := boardPermStat(user, uid, board, bid)
	if (board.Brdname[0] == '\x00') ||
		(board.BrdAttr&(ptttype.BRD_GROUPBOARD|ptttype.BRD_SYMBOLIC) != 0) ||
		!((state != ptttype.NBRD_INVALID) || isGroupOp) {
		return nil
	}

	boardStat := newBoardStat(bidInCache, state, board, isGroupOp)

	return boardStat
}

func LoadBoardsByBids(user *ptttype.UserecRaw, uid ptttype.Uid, bids []ptttype.Bid) (summaries []*ptttype.BoardSummaryRaw, err error) {

	boardStats := make([]*ptttype.BoardStat, 0, len(bids))

	for _, bid := range bids {
		if !bid.IsValid() {
			continue
		}
		eachBoardStat := loadBoardStat(user, uid, bid)
		if eachBoardStat == nil {
			continue
		}

		boardStats = append(boardStats, eachBoardStat)
	}

	summaries, err = showBoardList(user, uid, boardStats)

	return summaries, err
}

func loadBoardStat(user *ptttype.UserecRaw, uid ptttype.Uid, bid ptttype.Bid) (boardStat *ptttype.BoardStat) {
	bidInCache := bid.ToBidInStore()
	board := &ptttype.BoardHeaderRaw{}
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache),
		ptttype.BOARD_HEADER_RAW_SZ,
		unsafe.Pointer(board),
	)

	isGroupOp := groupOp(user, board)
	state := boardPermStat(user, uid, board, bid)
	if (board.Brdname[0] == '\x00') ||
		(board.BrdAttr&(ptttype.BRD_GROUPBOARD|ptttype.BRD_SYMBOLIC) != 0) ||
		!((state != ptttype.NBRD_INVALID) || isGroupOp) {
		return nil
	}

	boardStat = newBoardStat(bidInCache, state, board, isGroupOp)
	return boardStat
}

//LoadAutoCompleteBoards
//
//Load auto-complete boards by name.
func LoadAutoCompleteBoards(user *ptttype.UserecRaw, uid ptttype.Uid, startIdx ptttype.SortIdx, nBoards int, keyword []byte, isAsc bool) (summaries []*ptttype.BoardSummaryRaw, nextSummary *ptttype.BoardSummaryRaw, err error) {

	nBoardsInCache := cache.NumBoards()
	if startIdx == 0 && !isAsc {
		startIdx = ptttype.SortIdx(nBoardsInCache)
	}

	startIdxInStore := startIdx.ToSortIdxInStore()

	nBoardsWithNext := nBoards + 1

	//get board-stats
	boardStats := make([]*ptttype.BoardStat, 0, nBoardsWithNext)
	if isAsc {
		for idxInStore := startIdxInStore; ; idxInStore++ {
			if int32(idxInStore) >= nBoardsInCache || len(boardStats) >= nBoardsWithNext { //add 1 more board for nextSummary
				break
			}
			eachBoardStat, isEnd := loadAutoCompleteBoardStat(user, uid, idxInStore, keyword)
			if isEnd {
				break
			}
			if eachBoardStat == nil {
				continue
			}

			boardStats = append(boardStats, eachBoardStat)
		}
	} else {
		for idxInStore := startIdxInStore; ; idxInStore-- {
			if int32(idxInStore) < 0 || len(boardStats) >= nBoardsWithNext { //add 1 more board for nextSummary
				break
			}
			eachBoardStat, isEnd := loadAutoCompleteBoardStat(user, uid, idxInStore, keyword)
			if isEnd {
				break
			}
			if eachBoardStat == nil {
				continue
			}

			boardStats = append(boardStats, eachBoardStat)
		}
	}

	//boardStats to summaries
	summaries, err = showBoardList(user, uid, boardStats)
	if err != nil {
		return nil, nil, err
	}

	if len(summaries) == nBoardsWithNext {
		nextSummary = summaries[nBoards]
		summaries = summaries[:nBoards]
	}

	return summaries, nextSummary, nil
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
func LoadGeneralBoards(user *ptttype.UserecRaw, uid ptttype.Uid, startIdx ptttype.SortIdx, nBoards int, title []byte, keyword []byte, isAsc bool, bsortBy ptttype.BSortBy) (summaries []*ptttype.BoardSummaryRaw, nextSummary *ptttype.BoardSummaryRaw, err error) {

	nBoardsInCache := cache.NumBoards()
	if startIdx == 0 && !isAsc {
		startIdx = ptttype.SortIdx(nBoardsInCache)
	}

	startIdxInStore := startIdx.ToSortIdxInStore()

	nBoardsWithNext := nBoards + 1

	//get board-stats
	boardStats := make([]*ptttype.BoardStat, 0, nBoardsWithNext)
	if isAsc {
		for idxInStore := startIdxInStore; ; idxInStore++ {
			if int32(idxInStore) >= nBoardsInCache || len(boardStats) >= nBoardsWithNext { //add 1 more board for nextSummary
				break
			}
			eachBoardStat := loadGeneralBoardStat(user, uid, idxInStore, title, keyword, bsortBy)
			if eachBoardStat == nil {
				continue
			}

			boardStats = append(boardStats, eachBoardStat)
		}
	} else {
		for idxInStore := startIdxInStore; ; idxInStore-- {
			if int32(idxInStore) < 0 || len(boardStats) >= nBoardsWithNext { //add 1 more board for nextSummary
				break
			}
			eachBoardStat := loadGeneralBoardStat(user, uid, idxInStore, title, keyword, bsortBy)
			if eachBoardStat == nil {
				continue
			}

			boardStats = append(boardStats, eachBoardStat)
		}
	}

	//boardStats to summaries
	summaries, err = showBoardList(user, uid, boardStats)
	if err != nil {
		return nil, nil, err
	}

	if len(summaries) == nBoardsWithNext {
		nextSummary = summaries[nBoards]
		summaries = summaries[:nBoards]
	}

	return summaries, nextSummary, nil
}

//loadAutoCompleteBoardStat
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1147
func loadAutoCompleteBoardStat(user *ptttype.UserecRaw, uid ptttype.Uid, idxInStore ptttype.SortIdxInStore, keyword []byte) (boardStat *ptttype.BoardStat, isEnd bool) {
	var bidInCache ptttype.BidInStore

	const bsort0sz = unsafe.Sizeof(cache.Shm.Raw.BSorted[0])
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BSorted)+bsort0sz*uintptr(ptttype.BSORT_BY_NAME)+uintptr(idxInStore)*ptttype.BID_IN_STORE_SZ,
		ptttype.BID_IN_STORE_SZ,
		unsafe.Pointer(&bidInCache),
	)
	if bidInCache < 0 {
		return nil, false
	}

	board := &ptttype.BoardHeaderRaw{}
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache),
		ptttype.BOARD_HEADER_RAW_SZ,
		unsafe.Pointer(board),
	)

	if !types.CstrCaseHasPrefix(board.Brdname[:], keyword) {
		return nil, true
	}

	bid := bidInCache.ToBid()
	isGroupOp := groupOp(user, board)
	state := boardPermStat(user, uid, board, bid)
	if (board.Brdname[0] == '\x00') ||
		(board.BrdAttr&(ptttype.BRD_GROUPBOARD|ptttype.BRD_SYMBOLIC) != 0) ||
		!((state != ptttype.NBRD_INVALID) || isGroupOp) {
		return nil, false
	}

	boardStat = newBoardStat(bidInCache, state, board, isGroupOp)
	return boardStat, false
}

//loadGeneralBoardStat
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1147
func loadGeneralBoardStat(user *ptttype.UserecRaw, uid ptttype.Uid, idxInStore ptttype.SortIdxInStore, title []byte, keyword []byte, bsortBy ptttype.BSortBy) (boardStat *ptttype.BoardStat) {
	var bidInCache ptttype.BidInStore

	const bsort0sz = unsafe.Sizeof(cache.Shm.Raw.BSorted[0])
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BSorted)+bsort0sz*uintptr(bsortBy)+uintptr(idxInStore)*ptttype.BID_IN_STORE_SZ,
		ptttype.BID_IN_STORE_SZ,
		unsafe.Pointer(&bidInCache),
	)
	if bidInCache < 0 {
		return nil
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
		keywordsNotInBoard(&board.Brdname, &board.Title, title, keyword) {
		return nil
	}

	boardStat = newBoardStat(bidInCache, state, board, isGroupOp)
	return boardStat
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

//keywordsNotInBoard
//
//TITLE_MATCH in board.c
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L14
func keywordsNotInBoard(boardID *ptttype.BoardID_t, boardTitle *ptttype.BoardTitle_t, title []byte, keyword []byte) bool {
	if len(title) > 0 {
		return (types.Cstrcasestr(boardTitle[:], title) < 0)
	}

	if len(keyword) > 0 {
		if types.Cstrcasestr(boardTitle[:], keyword) >= 0 {
			return false
		}
		if types.Cstrcasestr(boardID[:], keyword) >= 0 {
			return false
		}

		return true
	}

	return false
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
	if !boardStat.IsGroupOp && boardStat.Attr == ptttype.NBRD_INVALID {
		summary = ptttype.NewBoardSummaryRawWithReason(boardStat)
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

	summary = ptttype.NewBoardSummaryRaw(boardStat, lastPostTime, total)

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

func FindBoardStartIdxByName(boardID *ptttype.BoardID_t, isAsc bool) (startIdx ptttype.SortIdx, err error) {

	return cache.FindBoardIdxByName(boardID, isAsc)
}

func FindBoardStartIdxByClass(cls []byte, boardID *ptttype.BoardID_t, isAsc bool) (startIdx ptttype.SortIdx, err error) {

	return cache.FindBoardIdxByClass(cls, boardID, isAsc)
}

func FindBoardAutoCompleteStartIdx(keyword []byte, isAsc bool) (startIdx ptttype.SortIdx, err error) {
	return cache.FindBoardAutoCompleteStartIdx(keyword, isAsc)
}
