package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// LoadFullClassBoards
//
// Load full class boards
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1142
//
// Params
//
//	user
//	uid
//	startIdx: the idx in bsorted.
//	nBoards: try to get at most nBoards
//
// Return
//
//	summaries
//	nextIdx: next idx in bsorted.
//	err
func LoadFullClassBoards(user *ptttype.UserecRaw, uid ptttype.UID, startBid ptttype.Bid, nBoards int) (summaries []*ptttype.BoardSummaryRaw, nextSummary *ptttype.BoardSummaryRaw, err error) {
	if !startBid.IsValid() {
		return nil, nil, ptttype.ErrInvalidBid
	}

	nBoardsInCache := cache.NumBoards()

	nBoardsWithNext := nBoards + 1

	// get board-stats
	boardStats := make([]*ptttype.BoardStat, 0, nBoardsWithNext)
	for bid := startBid; ; bid++ {
		bidInCache := bid.ToBidInStore()
		if int32(bidInCache) >= nBoardsInCache || len(boardStats) >= nBoardsWithNext { // add 1 more board for nextSummary
			break
		}
		eachBoardStat, _, eachErr := loadClassBoardStat(user, uid, bid, true)
		if eachErr != nil {
			continue
		}

		boardStats = append(boardStats, eachBoardStat)
	}

	// boardStats to summaries
	summaries, err = showBoardList(user, uid, boardStats, true)
	if err != nil {
		return nil, nil, err
	}

	if len(summaries) == nBoardsWithNext {
		nextSummary = summaries[nBoards]
		summaries = summaries[:nBoards]
	}

	return summaries, nextSummary, nil
}

// LoadClassBoards
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1169
func LoadClassBoards(user *ptttype.UserecRaw, uid ptttype.UID, classBid ptttype.Bid, bsortBy ptttype.BSortBy) (summaries []*ptttype.BoardSummaryRaw, err error) {
	if !classBid.IsValid() {
		return nil, ptttype.ErrInvalidBid
	}

	// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1041
	// override type in class root, because usually we don't need to sort
	// class root; and there may be out-of-sync in that mode.
	if isInClassRoot(classBid) {
		bsortBy = ptttype.BSORT_BY_CLASS
	}

	board, err := cache.GetBCache(classBid)
	if err != nil {
		return nil, err
	}

	if board.FirstChild[bsortBy] == 0 || board.ChildCount == 0 {
		err = cache.ResolveBoardGroup(classBid, bsortBy)
		if err != nil {
			return nil, err
		}

		board, err = cache.GetBCache(classBid)
		if err != nil {
			return nil, err
		}
	}

	// Ptt: child count after resolve_board_group
	childCount := int(board.ChildCount)

	brdSize := childCount + 5
	boardStats := make([]*ptttype.BoardStat, 0, brdSize)
	var eachBoardStat *ptttype.BoardStat
	var eachErr error
	for bid := board.FirstChild[bsortBy]; bid > 0 && len(boardStats) < brdSize; bid = board.Next[bsortBy] {
		if !bid.IsValid() {
			break
		}

		eachBoardStat, board, eachErr = loadClassBoardStat(user, uid, bid, false)
		if eachErr != nil {
			continue
		}
		boardStats = append(boardStats, eachBoardStat)
	}
	if childCount < len(boardStats) {
		// Ptt: dirty fix fix soon
		_ = cache.SetBoardChildCount(classBid, int32(0))
	}

	summaries, err = showBoardList(user, uid, boardStats, true)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func isInClassRoot(classBid ptttype.Bid) bool {
	return classBid == 1
}

func LoadBoardDetail(user *ptttype.UserecRaw, uid ptttype.UID, bid ptttype.Bid) (boardDetail *ptttype.BoardDetailRaw, err error) {
	if !bid.IsValid() {
		return nil, ptttype.ErrInvalidBid
	}

	bidInCache := bid.ToBidInStore()
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, err
	}

	isGroupOp := groupOp(user, uid, board)
	state := boardPermStat(user, uid, board, bid)
	boardStat := newBoardStat(bidInCache, state, board, isGroupOp)

	if boardStat == nil {
		return nil, err
	}

	lastPostTime, _ := cache.GetLastPosttime(boardStat.Bid)
	total, _ := cache.GetBTotalWithRetry(boardStat.Bid)

	boardDetail = &ptttype.BoardDetailRaw{
		Bid:            bid,
		LastPostTime:   lastPostTime,
		Total:          total,
		BoardHeaderRaw: board,
	}

	return boardDetail, nil
}

// LoadBoardSummary
//
// # Load General Board Summary
//
// Params:
//
//	user
//	uid
//	bid
//
// Return:
//
//	summary
//	err
func LoadBoardSummary(user *ptttype.UserecRaw, uid ptttype.UID, bid ptttype.Bid) (summary *ptttype.BoardSummaryRaw, err error) {
	if !bid.IsValid() {
		return nil, ptttype.ErrInvalidBid
	}

	bidInCache := bid.ToBidInStore()
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, err
	}
	isGroupOp := groupOp(user, uid, board)
	state := boardPermStat(user, uid, board, bid)
	boardStat := newBoardStat(bidInCache, state, board, isGroupOp)

	if boardStat == nil {
		return nil, err
	}

	isParseFolder := board.BrdAttr.HasPerm(ptttype.BRD_GROUPBOARD)
	summary = parseBoardSummary(user, uid, boardStat, isParseFolder)

	return summary, nil
}

// LoadHotBoards
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1125
func LoadHotBoards(user *ptttype.UserecRaw, uid ptttype.UID) (summary []*ptttype.BoardSummaryRaw, err error) {
	nBoards := cache.NHots()

	boardStats := make([]*ptttype.BoardStat, 0, nBoards)

	for idx := uint8(0); idx < nBoards; idx++ {
		eachBoardStat := loadHotBoardStat(user, uid, idx)

		if eachBoardStat == nil {
			continue
		}

		boardStats = append(boardStats, eachBoardStat)
	}

	summary, err = showBoardList(user, uid, boardStats, false)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// loadHotBoardStat
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1147
func loadHotBoardStat(user *ptttype.UserecRaw, uid ptttype.UID, idx uint8) *ptttype.BoardStat {
	// read bid-in-cache
	bidInCache := cache.Shm.Shm.HBcache[idx]
	if bidInCache < 0 {
		return nil
	}

	// get board
	board := &cache.Shm.Shm.BCache[bidInCache]

	// board-stat
	// assuming that the hot-boards can be accessed by the public.
	bid := bidInCache.ToBid()
	isGroupOp := groupOp(user, uid, board)
	state := boardPermStat(user, uid, board, bid)
	if (board.Brdname[0] == '\x00') ||
		(board.BrdAttr&(ptttype.BRD_GROUPBOARD|ptttype.BRD_SYMBOLIC) != 0) ||
		!((state != ptttype.NBRD_INVALID) || isGroupOp) {
		return nil
	}

	boardStat := newBoardStat(bidInCache, state, board, isGroupOp)

	return boardStat
}

func LoadBoardsByBids(user *ptttype.UserecRaw, uid ptttype.UID, bids []ptttype.Bid) (summaries []*ptttype.BoardSummaryRaw, err error) {
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

	summaries, err = showBoardList(user, uid, boardStats, true)

	return summaries, err
}

func loadBoardStat(user *ptttype.UserecRaw, uid ptttype.UID, bid ptttype.Bid) (boardStat *ptttype.BoardStat) {
	bidInCache := bid.ToBidInStore()
	board := &cache.Shm.Shm.BCache[bidInCache]

	isGroupOp := groupOp(user, uid, board)
	state := boardPermStat(user, uid, board, bid)
	if (board.Brdname[0] == '\x00') ||
		!((state != ptttype.NBRD_INVALID) || isGroupOp) {
		return nil
	}

	boardStat = newBoardStat(bidInCache, state, board, isGroupOp)
	return boardStat
}

// LoadAutoCompleteBoards
//
// Load auto-complete boards by name.
func LoadAutoCompleteBoards(user *ptttype.UserecRaw, uid ptttype.UID, startIdx ptttype.SortIdx, nBoards int, keyword []byte, isAsc bool) (summaries []*ptttype.BoardSummaryRaw, nextSummary *ptttype.BoardSummaryRaw, err error) {
	nBoardsInCache := cache.NumBoards()
	if startIdx == 0 && !isAsc {
		startIdx = ptttype.SortIdx(nBoardsInCache)
	}

	startIdxInStore := startIdx.ToSortIdxInStore()

	nBoardsWithNext := nBoards + 1

	// get board-stats
	boardStats := make([]*ptttype.BoardStat, 0, nBoardsWithNext)
	if isAsc {
		for idxInStore := startIdxInStore; ; idxInStore++ {
			if int32(idxInStore) >= nBoardsInCache || len(boardStats) >= nBoardsWithNext { // add 1 more board for nextSummary
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
			if int32(idxInStore) < 0 || len(boardStats) >= nBoardsWithNext { // add 1 more board for nextSummary
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

	// boardStats to summaries
	summaries, err = showBoardList(user, uid, boardStats, false)
	if err != nil {
		return nil, nil, err
	}

	if len(summaries) == nBoardsWithNext {
		nextSummary = summaries[nBoards]
		summaries = summaries[:nBoards]
	}

	return summaries, nextSummary, nil
}

// LoadGeneralBoards
//
// Load general boards by name.
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1142
// Params:
//
//		user
//		uid
//		startIdx: the idx in bsorted.
//	 nBoards: try to get at most nBoards
//		keyword
//
// Return:
//
//	 summary
//		nextIdx: next idx in bsorted.
//		err
func LoadGeneralBoards(user *ptttype.UserecRaw, uid ptttype.UID, startIdx ptttype.SortIdx, nBoards int, title []byte, keyword []byte, isAsc bool, bsortBy ptttype.BSortBy) (summaries []*ptttype.BoardSummaryRaw, nextSummary *ptttype.BoardSummaryRaw, err error) {
	nBoardsInCache := cache.NumBoards()
	if startIdx == 0 && !isAsc {
		startIdx = ptttype.SortIdx(nBoardsInCache)
	}

	startIdxInStore := startIdx.ToSortIdxInStore()

	nBoardsWithNext := nBoards + 1

	// get board-stats
	boardStats := make([]*ptttype.BoardStat, 0, nBoardsWithNext)
	if isAsc {
		for idxInStore := startIdxInStore; ; idxInStore++ {
			if int32(idxInStore) >= nBoardsInCache || len(boardStats) >= nBoardsWithNext { // add 1 more board for nextSummary
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
			if int32(idxInStore) < 0 || len(boardStats) >= nBoardsWithNext { // add 1 more board for nextSummary
				break
			}
			eachBoardStat := loadGeneralBoardStat(user, uid, idxInStore, title, keyword, bsortBy)
			if eachBoardStat == nil {
				continue
			}

			boardStats = append(boardStats, eachBoardStat)
		}
	}

	// boardStats to summaries
	summaries, err = showBoardList(user, uid, boardStats, false)
	if err != nil {
		return nil, nil, err
	}

	if len(summaries) == nBoardsWithNext {
		nextSummary = summaries[nBoards]
		summaries = summaries[:nBoards]
	}

	return summaries, nextSummary, nil
}

// loadAutoCompleteBoardStat
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1147
func loadAutoCompleteBoardStat(user *ptttype.UserecRaw, uid ptttype.UID, idxInStore ptttype.SortIdxInStore, keyword []byte) (boardStat *ptttype.BoardStat, isEnd bool) {
	bidInCache := cache.Shm.Shm.BSorted[ptttype.BSORT_BY_NAME][idxInStore]
	if bidInCache < 0 {
		return nil, false
	}

	board := &cache.Shm.Shm.BCache[bidInCache]

	if !types.CstrCaseHasPrefix(board.Brdname[:], keyword) {
		return nil, true
	}

	bid := bidInCache.ToBid()
	isGroupOp := groupOp(user, uid, board)
	state := boardPermStat(user, uid, board, bid)
	if (board.Brdname[0] == '\x00') ||
		(board.BrdAttr&(ptttype.BRD_GROUPBOARD|ptttype.BRD_SYMBOLIC) != 0) ||
		!((state != ptttype.NBRD_INVALID) || isGroupOp) {
		return nil, false
	}

	boardStat = newBoardStat(bidInCache, state, board, isGroupOp)
	return boardStat, false
}

// loadClassBoardStat
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1186
func loadClassBoardStat(user *ptttype.UserecRaw, uid ptttype.UID, bid ptttype.Bid, isResolveBoardGroup bool) (boardStat *ptttype.BoardStat, board *ptttype.BoardHeaderRaw, err error) {
	bsortBy := ptttype.BSORT_BY_CLASS

	board, err = cache.GetBCache(bid)
	if err != nil {
		return nil, nil, err
	}

	if board.Brdname[0] == '\x00' ||
		!board.BrdAttr.HasPerm(ptttype.BRD_GROUPBOARD|ptttype.BRD_SYMBOLIC) {
		return nil, nil, ErrInvalidBoard
	}

	if isResolveBoardGroup && (board.FirstChild[bsortBy] == 0 || board.ChildCount == 0) {
		err = cache.ResolveBoardGroup(bid, bsortBy)
		if err != nil {
			return nil, nil, err
		}

		board, err = cache.GetBCache(bid)
		if err != nil {
			return nil, nil, err
		}
	}

	isGroupOp := groupOp(user, uid, board)
	state := boardPermStat(user, uid, board, bid)
	if !((state != ptttype.NBRD_INVALID) || isGroupOp) {
		return nil, nil, ErrNotPermitted
	}

	bidInCache := bid.ToBidInStore()
	boardStat = newBoardStat(bidInCache, state, board, isGroupOp)
	return boardStat, board, nil
}

// loadGeneralBoardStat
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1147
func loadGeneralBoardStat(user *ptttype.UserecRaw, uid ptttype.UID, idxInStore ptttype.SortIdxInStore, title []byte, keyword []byte, bsortBy ptttype.BSortBy) (boardStat *ptttype.BoardStat) {
	bidInCache := cache.Shm.Shm.BSorted[bsortBy][idxInStore]
	if bidInCache < 0 {
		return nil
	}

	board := &cache.Shm.Shm.BCache[bidInCache]

	bid := bidInCache.ToBid()
	isGroupOp := groupOp(user, uid, board)
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

// newBoardStat
func newBoardStat(bidInCache ptttype.BidInStore, state ptttype.BoardStatAttr, board *ptttype.BoardHeaderRaw, isGroupOp bool) (boardStat *ptttype.BoardStat) {
	boardStat = &ptttype.BoardStat{}

	boardStat.Bid = bidInCache.ToBid()
	boardStat.Attr = state

	boardStat.Board = board
	boardStat.IsGroupOp = isGroupOp

	// XXX need to modify this by having state with NBRD_SET_POSTMASK
	// XXX this is a hack to ensure the brd-postmask
	brd_postmask := ptttype.BRD_POSTMASK
	if (board.BrdAttr&ptttype.BRD_HIDE != 0) && (board.BrdAttr&ptttype.BRD_POSTMASK == 0) && state == ptttype.NBRD_BOARD {
		cache.Shm.Shm.BCache[bidInCache].BrdAttr |= brd_postmask
		board.BrdAttr |= brd_postmask
	}

	return boardStat
}

// keywordsNotInBoard
//
// TITLE_MATCH in board.c
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L14
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

// showBoardList
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1409
func showBoardList(user *ptttype.UserecRaw, uid ptttype.UID, boardStats []*ptttype.BoardStat, isParseFolder bool) (summary []*ptttype.BoardSummaryRaw, err error) {
	summary = make([]*ptttype.BoardSummaryRaw, len(boardStats))
	for idx, eachStat := range boardStats {
		summary[idx] = parseBoardSummary(user, uid, eachStat, isParseFolder)
	}

	return summary, nil
}

// parseBoardSummary
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1460
func parseBoardSummary(user *ptttype.UserecRaw, uid ptttype.UID, boardStat *ptttype.BoardStat, isParseFolder bool) (summary *ptttype.BoardSummaryRaw) {
	// XXX we do not deal with fav in go-pttbbs.
	if boardStat.Attr&ptttype.NBRD_LINE != 0 {
		return &ptttype.BoardSummaryRaw{Bid: boardStat.Bid, StatAttr: boardStat.Attr}
	}

	// XXX we do not deal with fav in go-pttbbs.
	if !isParseFolder && boardStat.Attr&ptttype.NBRD_FOLDER != 0 {
		return &ptttype.BoardSummaryRaw{Bid: boardStat.Bid, StatAttr: boardStat.Attr}
	}

	// hidden board
	if !boardStat.IsGroupOp && boardStat.Attr == ptttype.NBRD_INVALID {
		summary = ptttype.NewBoardSummaryRawWithReason(boardStat)
		return summary
	}

	lastPostTime, _ := cache.GetLastPosttime(boardStat.Bid)
	total, _ := cache.GetBTotalWithRetry(boardStat.Bid)

	summary = ptttype.NewBoardSummaryRaw(boardStat, lastPostTime, total)

	return summary
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
