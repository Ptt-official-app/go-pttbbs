package cache

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func GetBCache(bid ptttype.Bid) (board *ptttype.BoardHeaderRaw, err error) {
	if !bid.IsValid() {
		return nil, ptttype.ErrInvalidBid
	}

	bidInCache := bid.ToBidInStore()

	board = &ptttype.BoardHeaderRaw{}

	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache),
		ptttype.BOARD_HEADER_RAW_SZ,
		unsafe.Pointer(board),
	)
	return board, nil
}

func GetBTotalWithRetry(bid ptttype.Bid) (total int32, err error) {
	if !bid.IsValid() {
		return 0, ptttype.ErrInvalidBid
	}

	// 2. bcache preparation.
	total = GetBTotal(bid)
	if total == 0 {
		err = SetBTotal(bid)
		if err != nil {
			return 0, err
		}
		err = SetBottomTotal(bid)
		if err != nil {
			return 0, err
		}

		total = GetBTotal(bid)
		if total == 0 { // no data
			return 0, nil
		}
	}

	return total, nil
}

func GetBottomTotal(bid ptttype.Bid) (total int32) {
	if !bid.IsValid() {
		return 0
	}

	bidInCache := bid.ToBidInStore()

	var total_u8 uint8
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.NBottom)+types.UINT8_SZ*uintptr(bidInCache),
		types.UINT8_SZ,
		unsafe.Pointer(&total_u8),
	)

	return int32(total_u8)
}

func GetBTotal(bid ptttype.Bid) (total int32) {
	if !bid.IsValid() {
		return 0
	}

	bidInCache := bid.ToBidInStore()

	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Total)+types.INT32_SZ*uintptr(bidInCache),
		types.INT32_SZ,
		unsafe.Pointer(&total),
	)

	return total
}

//SetBTotal
//
//It's possible that we loaded nothing from ReloadBCache in the beginning of the program, and then there are some articles after a while.
//We need to sync the btotal and lastposttime back to shm.
func SetBTotal(bid ptttype.Bid) (err error) {
	if !bid.IsValid() {
		return ptttype.ErrInvalidBid
	}

	board, err := GetBCache(bid)
	if err != nil {
		return err
	}
	dirFilename, err := path.SetBFile(&board.Brdname, ptttype.FN_DIR)
	if err != nil {
		return err
	}

	file, err := os.Open(dirFilename)
	if err != nil { // we should always have .DIR
		if os.IsNotExist(err) {
			err = nil
		}
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	nArticles := int32(stat.Size() / int64(ptttype.FILE_HEADER_RAW_SZ))

	bidInCache := bid.ToBidInStore()
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Total)+types.INT32_SZ*uintptr(bidInCache),
		types.INT32_SZ,
		unsafe.Pointer(&nArticles),
	)

	zero := types.Time4(0)
	if nArticles == 0 {
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.LastPostTime)+types.TIME4_SZ*uintptr(bidInCache),
			types.TIME4_SZ,
			unsafe.Pointer(&zero),
		)
		return nil
	}

	// https://github.com/ptt/pttbbs/blob/master/common/bbs/cache.c#L633
	// lastPostTime is in filename (create-time), starting from 2nd bytes
	// ptttype.FileHeaderRaw
	_, err = file.Seek(int64(nArticles-1)*int64(ptttype.FILE_HEADER_RAW_SZ), 0)
	if err != nil {
		return err
	}

	articleFilename := &ptttype.Filename_t{}
	err = types.BinaryRead(file, binary.LittleEndian, articleFilename)
	if err != nil {
		return err
	}
	if types.Cstrcmp(articleFilename[:], []byte(ptttype.FN_SAFEDEL)) == 0 {
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.LastPostTime)+types.TIME4_SZ*uintptr(bidInCache),
			types.TIME4_SZ,
			unsafe.Pointer(&zero),
		)
		return nil
	}

	createTime, err := articleFilename.CreateTime()
	if err != nil {
		return err
	}

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.LastPostTime)+types.TIME4_SZ*uintptr(bidInCache),
		types.TIME4_SZ,
		unsafe.Pointer(&createTime),
	)

	return nil
}

func SetBottomTotal(bid ptttype.Bid) error {
	if !bid.IsValid() {
		return ptttype.ErrInvalidBid
	}

	board, err := GetBCache(bid)
	if err != nil {
		return err
	}
	if board.Brdname[0] == 0 {
		return nil
	}

	bottomFilename, err := path.SetBFile(&board.Brdname, ptttype.FN_DIR_BOTTOM)
	if err != nil {
		return err
	}

	bidInCache := bid.ToBidInStore()
	zero8 := uint8(0)
	n := uint8(cmsys.GetNumRecords(bottomFilename, ptttype.FILE_HEADER_RAW_SZ))
	if n > 5 {
		_ = syscall.Unlink(bottomFilename)
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.NBottom)+types.UINT8_SZ*uintptr(bidInCache),
			types.UINT8_SZ,
			unsafe.Pointer(&zero8),
		)

		return nil
	}

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.NBottom)+types.UINT8_SZ*uintptr(bidInCache),
		types.UINT8_SZ,
		unsafe.Pointer(&n),
	)

	return nil
}

func IsHiddenBoardFriend(bidInCache ptttype.BidInStore, uidInCache ptttype.UIDInStore) bool {
	if !bidInCache.ToBid().IsValid() || !uidInCache.ToUID().IsValid() {
		return false
	}

	// hbfl time
	var loadTime types.Time4
	pLoadTime := &loadTime

	const Hbfl0Size = unsafe.Sizeof(Shm.Raw.Hbfl[0])
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Hbfl)+Hbfl0Size*uintptr(bidInCache),
		types.TIME4_SZ,
		unsafe.Pointer(pLoadTime),
	)

	// XXX use nowTS to replace loginStartTime.
	//     HBFLexpire is set as 5-days. nowTS should be ok.
	nowTS := types.NowTS()
	if loadTime < nowTS-types.Time4(ptttype.HBFLexpire) {
		HbflReload(bidInCache)
	}

	uid := uidInCache.ToUID()

	var friendID ptttype.UID
	pFriendID := &friendID
	friendIDptr := unsafe.Pointer(pFriendID)
	for i := uintptr(1); i <= ptttype.MAX_FRIEND; i++ {
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.Hbfl)+Hbfl0Size*uintptr(bidInCache)+types.INT32_SZ*i,
			types.INT32_SZ,
			friendIDptr,
		)
		if friendID == 0 {
			break
		}

		if friendID == uid {
			return true
		}
	}

	return false
}

func HbflReload(bidInCache ptttype.BidInStore) {
	if !bidInCache.ToBid().IsValid() {
		return
	}

	brdname := &ptttype.BoardID_t{}

	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache)+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
		ptttype.BOARD_ID_SZ,
		unsafe.Pointer(brdname),
	)

	filename, err := path.SetBFile(brdname, ptttype.FN_VISIBLE)
	if err != nil {
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	hbfl := [ptttype.MAX_FRIEND + 1]ptttype.UID{}
	const hbflsz = unsafe.Sizeof(hbfl)

	reader := bufio.NewReader(file)
	var line []byte
	var uid ptttype.UID
	// num++ is in the end of the for.
	for num := ptttype.UID(1); num <= ptttype.MAX_FRIEND; {
		line, _ = types.ReadLine(reader)
		if len(line) == 0 {
			break
		}
		theList := bytes.Split(line, []byte{' '}) // The \x00 is taken care of by scanner.

		eachUserID := &ptttype.UserID_t{}
		copy(eachUserID[:], theList[0][:])

		if types.Cstrcasecmp(eachUserID[:], ptttype.USER_ID_GUEST[:]) == 0 {
			continue
		}

		uid, err = SearchUserRaw(eachUserID, nil)
		if err != nil {
			continue
		}
		if uid == 0 {
			continue
		}

		hbfl[num] = uid

		num++ // num++ is in the end of the for. (no num++ for the continue conditions)
	}

	hbfl[0] = ptttype.UID(types.NowTS())

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Hbfl)+hbflsz*uintptr(bidInCache),
		hbflsz,
		unsafe.Pointer(&hbfl),
	)
}

//NumBoards
//
//https://github.com/ptt/pttbbs/blob/master/common/bbs/cache.c#L512
func NumBoards() int32 {
	return Shm.GetBNumber()
}

func NHots() (nhots uint8) {
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.NHOTs),
		types.UINT8_SZ,
		unsafe.Pointer(&nhots),
	)

	return nhots
}

//Reload BCache
//
//https://github.com/ptt/pttbbs/blob/master/common/bbs/cache.c#L458
func ReloadBCache() {
	var busystate int32
	for i := 0; i < 10; i++ { // Is it ok that we don't use mutex or semaphore here?
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.BBusyState),
			types.INT32_SZ,
			unsafe.Pointer(&busystate),
		)
		if busystate == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	// XXX should we check that the busystate is still != 0 and return?
	reloadBCacheCore()

	SortBCache()

	// XXX reloadCacheLoadBottom as front-process
	reloadCacheLoadBottom()
}

func reloadBCacheCore() {
	busystate := int32(1)
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BBusyState),
		types.INT32_SZ,
		unsafe.Pointer(&busystate),
	)
	defer func() {
		busystate = 0
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.BBusyState),
			types.INT32_SZ,
			unsafe.Pointer(&busystate),
		)
	}()

	theBytes, err := reloadBCacheReadFile()
	if err != nil {
		return
	}

	const bcachesz = unsafe.Sizeof(Shm.Raw.BCache)
	theSize := bcachesz
	lenTheBytes := uintptr(len(theBytes))
	if lenTheBytes < theSize {
		theSize = lenTheBytes
	}
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BCache),
		theSize,
		unsafe.Pointer(&theBytes[0]),
	)

	bnumber := int32(theSize / ptttype.BOARD_HEADER_RAW_SZ)
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BNumber),
		types.INT32_SZ,
		unsafe.Pointer(&bnumber),
	)

	Shm.Memset(
		unsafe.Offsetof(Shm.Raw.LastPostTime),
		byte(0),
		uintptr(ptttype.MAX_BOARD)*types.TIME4_SZ,
	)

	Shm.Memset(
		unsafe.Offsetof(Shm.Raw.Total),
		byte(0),
		uintptr(ptttype.MAX_BOARD)*types.INT32_SZ,
	)

	Shm.InnerSetInt32(
		unsafe.Offsetof(Shm.Raw.BTouchTime),
		unsafe.Offsetof(Shm.Raw.BUptime),
	)
}

// SortBCache
// XXX TODO: implement
func SortBCache() {
	var busystate int32
	pbusystate := &busystate
	pbusystateptr := unsafe.Pointer(pbusystate)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BBusyState),
		types.INT32_SZ,
		pbusystateptr,
	)

	if busystate != 0 {
		time.Sleep(1 * time.Second)
		return
	}

	*pbusystate = 1
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BBusyState),
		types.INT32_SZ,
		pbusystateptr,
	)
	defer func() {
		*pbusystate = 0
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.BBusyState),
			types.INT32_SZ,
			pbusystateptr,
		)
	}()

	// init is in shm.c
	Shm.QsortCmpBoardName()
	Shm.QsortCmpBoardClass()

	// for-loop cleaning first-child
	// init vars
	val := [ptttype.BSORT_BY_MAX]ptttype.Bid{}
	bnumber := Shm.GetBNumber()

	const valsz = unsafe.Sizeof(val)
	for i := int32(0); i < bnumber; i++ {
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(i)+ptttype.BOARD_HEADER_FIRST_CHILD_OFFSET,
			valsz,
			unsafe.Pointer(&val),
		)
	}
}

func reloadCacheLoadBottom() {
	boardName := &ptttype.BoardID_t{}
	for i := uintptr(0); i < ptttype.MAX_BOARD; i++ {
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*i+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
			ptttype.BOARD_HEADER_RAW_SZ,
			unsafe.Pointer(boardName),
		)

		if boardName[0] == 0 {
			continue
		}

		filename, err := path.SetBFile(boardName, ptttype.FN_DIR_BOTTOM)
		if err != nil {
			continue
		}

		n := cmsys.GetNumRecords(filename, ptttype.FILE_HEADER_RAW_SZ)
		if n > 5 {
			n = 5
		}

		n_uint8 := uint8(n)
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.NBottom)+types.UINT8_SZ*i,
			types.UINT8_SZ,
			unsafe.Pointer(&n_uint8),
		)
	}
}

func reloadBCacheReadFile() ([]byte, error) {
	file, err := os.Open(ptttype.FN_BOARD)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	theBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return theBytes, nil
}

func GetBid(boardID *ptttype.BoardID_t) (bid ptttype.Bid, err error) {
	_, bid, err = getBidByNameCore(boardID)
	if err == ErrNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return bid, nil
}

func getBidByNameCore(boardID *ptttype.BoardID_t) (idx ptttype.SortIdxInStore, bid ptttype.Bid, err error) {
	// wait 1 second for bbusystate
	bbusystate := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BBusyState),
		types.INT32_SZ,
		unsafe.Pointer(&bbusystate),
	)
	if bbusystate != 0 {
		time.Sleep(1 * time.Second)
	}

	// start and end
	start := int32(0)
	end := Shm.GetBNumber()
	end--
	if end < 0 { // unable to get bid
		return -1, 0, nil
	}

	const bsort0sz = unsafe.Sizeof(Shm.Raw.BSorted[0])
	const bsortedOffset = unsafe.Offsetof(Shm.Raw.BSorted) + bsort0sz*uintptr(ptttype.BSORT_BY_NAME)
	const bcacheOffset = unsafe.Offsetof(Shm.Raw.BCache)
	bidInCache := ptttype.BidInStore(0)
	bidInCache_ptr := unsafe.Pointer(&bidInCache)
	boardIDInCache := &ptttype.BoardID_t{}
	boardIDInCache_ptr := unsafe.Pointer(boardIDInCache)
	idx_i32 := (start + end) / 2
	for ; ; idx_i32 = (start + end) / 2 {
		Shm.ReadAt(
			bsortedOffset+uintptr(idx_i32)*ptttype.BID_IN_STORE_SZ,
			ptttype.BID_IN_STORE_SZ,
			bidInCache_ptr,
		)

		Shm.ReadAt(
			bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
			ptttype.BOARD_ID_SZ,
			boardIDInCache_ptr,
		)

		j := types.Cstrcasecmp(boardID[:], boardIDInCache[:])
		if j == 0 {
			bid = bidInCache.ToBid()
			return ptttype.SortIdxInStore(idx_i32), bid, nil
		}

		if end == start {
			break
		} else if idx_i32 == start {
			idx_i32 = end // nolint
			start = end
		} else if j > 0 {
			start = idx_i32
		} else {
			end = idx_i32
		}
	}

	return ptttype.SortIdxInStore(idx_i32), 0, ErrNotFound
}

func getBidByClassCore(cls []byte, boardID *ptttype.BoardID_t) (idx ptttype.SortIdxInStore, bid ptttype.Bid, err error) {
	// wait 1 second for bbusystate
	bbusystate := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BBusyState),
		types.INT32_SZ,
		unsafe.Pointer(&bbusystate),
	)
	if bbusystate != 0 {
		time.Sleep(1 * time.Second)
	}

	// start and end
	start := int32(0)
	end := Shm.GetBNumber()
	end--
	if end < 0 { // unable to get bid
		return -1, 0, nil
	}

	const bsort0sz = unsafe.Sizeof(Shm.Raw.BSorted[0])
	const bsortedOffset = unsafe.Offsetof(Shm.Raw.BSorted) + bsort0sz*uintptr(ptttype.BSORT_BY_CLASS)
	const bcacheOffset = unsafe.Offsetof(Shm.Raw.BCache)
	bidInCache := ptttype.BidInStore(0)
	bidInCache_ptr := unsafe.Pointer(&bidInCache)

	titleInCache := &ptttype.BoardTitle_t{}
	titleInCache_ptr := unsafe.Pointer(titleInCache)

	boardIDInCache := &ptttype.BoardID_t{}
	boardIDInCache_ptr := unsafe.Pointer(boardIDInCache)
	idx_i32 := (start + end) / 2
	for ; ; idx_i32 = (start + end) / 2 {
		Shm.ReadAt(
			bsortedOffset+uintptr(idx_i32)*ptttype.BID_IN_STORE_SZ,
			ptttype.BID_IN_STORE_SZ,
			bidInCache_ptr,
		)

		Shm.ReadAt(
			bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_TITLE_OFFSET,
			ptttype.BOARD_TITLE_SZ,
			titleInCache_ptr,
		)
		clsInCache := titleInCache.BoardClass()

		Shm.ReadAt(
			bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
			ptttype.BOARD_ID_SZ,
			boardIDInCache_ptr,
		)

		j := cmpBoardByClass(cls, boardID, clsInCache, boardIDInCache)
		if j == 0 {
			bid = bidInCache.ToBid()
			return ptttype.SortIdxInStore(idx_i32), bid, nil
		}

		if end == start {
			break
		} else if idx_i32 == start {
			idx_i32 = end // nolint
			start = end
		} else if j > 0 {
			start = idx_i32
		} else {
			end = idx_i32
		}
	}

	return ptttype.SortIdxInStore(idx_i32), 0, ErrNotFound
}

func FindBoardIdxByName(boardID *ptttype.BoardID_t, isAsc bool) (idx ptttype.SortIdx, err error) {
	idxInStore, bid, err := getBidByNameCore(boardID)
	if bid.IsValid() && err == nil {
		return idxInStore.ToSortIdx(), nil
	}
	if err != ErrNotFound || idxInStore == -1 {
		return -1, err
	}

	nBoard_i32 := Shm.GetBNumber()
	nBoard := ptttype.SortIdxInStore(nBoard_i32)
	const bsort0sz = unsafe.Sizeof(Shm.Raw.BSorted[0])
	const bsortedOffset = unsafe.Offsetof(Shm.Raw.BSorted) + bsort0sz*uintptr(ptttype.BSORT_BY_NAME)
	const bcacheOffset = unsafe.Offsetof(Shm.Raw.BCache)
	bidInCache := ptttype.BidInStore(0)
	bidInCache_ptr := unsafe.Pointer(&bidInCache)
	boardIDInCache := &ptttype.BoardID_t{}
	boardIDInCache_ptr := unsafe.Pointer(boardIDInCache)
	//
	if isAsc {
		for ; idxInStore < nBoard; idxInStore++ {
			Shm.ReadAt(
				bsortedOffset+uintptr(idxInStore)*ptttype.BID_IN_STORE_SZ,
				ptttype.BID_IN_STORE_SZ,
				bidInCache_ptr,
			)

			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}

			Shm.ReadAt(
				bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
				ptttype.BOARD_ID_SZ,
				boardIDInCache_ptr,
			)

			j := types.Cstrcasecmp(boardID[:], boardIDInCache[:])
			if j < 0 {
				break
			}
		}
		if idxInStore == nBoard {
			idxInStore = -1
		}
	} else {
		for ; idxInStore >= 0; idxInStore-- {
			Shm.ReadAt(
				bsortedOffset+uintptr(idxInStore)*ptttype.BID_IN_STORE_SZ,
				ptttype.BID_IN_STORE_SZ,
				bidInCache_ptr,
			)

			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}

			Shm.ReadAt(
				bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
				ptttype.BOARD_ID_SZ,
				boardIDInCache_ptr,
			)

			j := types.Cstrcasecmp(boardID[:], boardIDInCache[:])
			if j > 0 {
				break
			}
		}
	}
	if idxInStore == -1 {
		return -1, nil
	}

	return idxInStore.ToSortIdx(), nil
}

func FindBoardIdxByClass(cls []byte, boardID *ptttype.BoardID_t, isAsc bool) (idx ptttype.SortIdx, err error) {
	idxInStore, bid, err := getBidByClassCore(cls, boardID)
	if bid.IsValid() && err == nil {
		return idxInStore.ToSortIdx(), nil
	}
	if err != ErrNotFound || idxInStore == -1 {
		return -1, err
	}

	nBoard_i32 := Shm.GetBNumber()
	nBoard := ptttype.SortIdxInStore(nBoard_i32)
	const bsort0sz = unsafe.Sizeof(Shm.Raw.BSorted[0])
	const bsortedOffset = unsafe.Offsetof(Shm.Raw.BSorted) + bsort0sz*uintptr(ptttype.BSORT_BY_CLASS)
	const bcacheOffset = unsafe.Offsetof(Shm.Raw.BCache)
	bidInCache := ptttype.BidInStore(0)
	bidInCache_ptr := unsafe.Pointer(&bidInCache)

	titleInCache := &ptttype.BoardTitle_t{}
	titleInCache_ptr := unsafe.Pointer(titleInCache)

	boardIDInCache := &ptttype.BoardID_t{}
	boardIDInCache_ptr := unsafe.Pointer(boardIDInCache)
	//
	if isAsc {
		for ; idxInStore < nBoard; idxInStore++ {
			Shm.ReadAt(
				bsortedOffset+uintptr(idxInStore)*ptttype.BID_IN_STORE_SZ,
				ptttype.BID_IN_STORE_SZ,
				bidInCache_ptr,
			)

			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}

			Shm.ReadAt(
				bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_TITLE_OFFSET,
				ptttype.BOARD_TITLE_SZ,
				titleInCache_ptr,
			)

			clsInCache := titleInCache.BoardClass()

			Shm.ReadAt(
				bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
				ptttype.BOARD_ID_SZ,
				boardIDInCache_ptr,
			)

			j := cmpBoardByClass(cls, boardID, clsInCache, boardIDInCache)
			if j < 0 {
				break
			}
		}
		if idxInStore == nBoard {
			idxInStore = -1
		}
	} else {
		for ; idxInStore >= 0; idxInStore-- {
			Shm.ReadAt(
				bsortedOffset+uintptr(idxInStore)*ptttype.BID_IN_STORE_SZ,
				ptttype.BID_IN_STORE_SZ,
				bidInCache_ptr,
			)

			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}

			Shm.ReadAt(
				bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_TITLE_OFFSET,
				ptttype.BOARD_TITLE_SZ,
				titleInCache_ptr,
			)
			clsInCache := titleInCache.BoardClass()

			Shm.ReadAt(
				bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
				ptttype.BOARD_ID_SZ,
				boardIDInCache_ptr,
			)

			j := cmpBoardByClass(cls, boardID, clsInCache, boardIDInCache)
			if j > 0 {
				break
			}
		}
	}
	if idxInStore == -1 {
		return -1, nil
	}

	return idxInStore.ToSortIdx(), nil
}

func cmpBoardByClass(cls []byte, boardID *ptttype.BoardID_t, clsInCache []byte, boardIDInCache *ptttype.BoardID_t) int {
	j := types.Cstrcmp(cls, clsInCache)
	if j != 0 {
		return j
	}

	return types.Cstrcasecmp(boardID[:], boardIDInCache[:])
}

func FindBoardAutoCompleteStartIdx(keyword []byte, isAsc bool) (startIdx ptttype.SortIdx, err error) {
	boardID := findBoardClosetKeyword(keyword, isAsc)
	nBoard_i32 := Shm.GetBNumber()
	nBoard := ptttype.SortIdxInStore(nBoard_i32)

	// find the closet keyword
	idx, err := FindBoardIdxByName(boardID, !isAsc)
	if err != nil {
		return -1, err
	}
	if idx == -1 {
		if isAsc {
			idx = 1
		} else {
			idx = ptttype.SortIdx(nBoard_i32)
		}
	}
	idxInStore := idx.ToSortIdxInStore()

	const bsort0sz = unsafe.Sizeof(Shm.Raw.BSorted[0])
	const bsortedOffset = unsafe.Offsetof(Shm.Raw.BSorted) + bsort0sz*uintptr(ptttype.BSORT_BY_NAME)
	const bcacheOffset = unsafe.Offsetof(Shm.Raw.BCache)
	bidInCache := ptttype.BidInStore(0)
	bidInCache_ptr := unsafe.Pointer(&bidInCache)

	boardIDInCache := &ptttype.BoardID_t{}
	boardIDInCache_ptr := unsafe.Pointer(boardIDInCache)

	// it should be either current idx or the next idx
	const MAX_ITER_FIND_AUTO_COMPLETE = 3
	if isAsc {
		i := 0
		for ; i < MAX_ITER_FIND_AUTO_COMPLETE && idxInStore < nBoard; i, idxInStore = i+1, idxInStore+1 {
			Shm.ReadAt(
				bsortedOffset+uintptr(idxInStore)*ptttype.BID_IN_STORE_SZ,
				ptttype.BID_IN_STORE_SZ,
				bidInCache_ptr,
			)

			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}

			Shm.ReadAt(
				bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
				ptttype.BOARD_ID_SZ,
				boardIDInCache_ptr,
			)

			boardIDInCachePrefix := boardIDInCache[:len(keyword)]
			j := types.Cstrcasecmp(keyword, boardIDInCachePrefix)
			if j == 0 {
				break
			} else if j < 0 { // keyword is already < boardIDInCachePrefix, we can't find fit anymore.
				return -1, nil
			}
		}
		if i == MAX_ITER_FIND_AUTO_COMPLETE || idxInStore == nBoard {
			idxInStore = -1
		}
	} else {
		i := 0
		for ; i < MAX_ITER_FIND_AUTO_COMPLETE && idxInStore >= 0; i, idxInStore = i+1, idxInStore-1 {
			Shm.ReadAt(
				bsortedOffset+uintptr(idxInStore)*ptttype.BID_IN_STORE_SZ,
				ptttype.BID_IN_STORE_SZ,
				bidInCache_ptr,
			)

			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}

			Shm.ReadAt(
				bcacheOffset+uintptr(bidInCache)*ptttype.BOARD_HEADER_RAW_SZ+ptttype.BOARD_HEADER_BRDNAME_OFFSET,
				ptttype.BOARD_ID_SZ,
				boardIDInCache_ptr,
			)

			boardIDInCachePrefix := boardIDInCache[:len(keyword)]
			j := types.Cstrcasecmp(keyword, boardIDInCachePrefix)
			if j == 0 {
				break
			} else if j > 0 { // keyword is already > boardIDInCachePrefix, we can't find fit anymore.
				return -1, nil
			}
		}
		if i == MAX_ITER_FIND_AUTO_COMPLETE {
			idxInStore = -1
		}
	}
	if idxInStore == -1 {
		return -1, nil
	}

	return idxInStore.ToSortIdx(), nil
}

func findBoardClosetKeyword(keyword []byte, isAsc bool) (boardID *ptttype.BoardID_t) {
	boardID = &ptttype.BoardID_t{}

	if isAsc {
		copy(boardID[:], keyword)
	} else {
		copy(boardID[:], keyword)
		boardID[len(keyword)-1]++
	}

	return boardID
}

func SanitizeBMs(bms *ptttype.BM_t) (parsedBMs *ptttype.BM_t) {
	if bms == nil {
		return &ptttype.BM_t{}
	}
	bmsBytes := types.CstrToBytes(bms[:])
	userIDsBytes := bytes.Split(bmsBytes, []byte{'/'})
	userIDs := make([]*ptttype.UserID_t, len(userIDsBytes))
	for idx, each := range userIDsBytes {
		userIDs[idx] = &ptttype.UserID_t{}
		copy(userIDs[idx][:], each)
	}

	validUserIDs := make([]*ptttype.UserID_t, 0, len(userIDs))
	for _, each := range userIDs {
		uid, err := SearchUserRaw(each, nil)
		if err != nil || uid == 0 {
			continue
		}
		validUserIDs = append(validUserIDs, each)
	}

	parsedBMs = ptttype.NewBM(validUserIDs)

	return parsedBMs
}

func ParseBMList(bms *ptttype.BM_t) (uids *[ptttype.MAX_BMs]ptttype.UID) {
	// init uids
	uids = &[ptttype.MAX_BMs]ptttype.UID{}
	for idx := 0; idx < ptttype.MAX_BMs; idx++ {
		uids[idx] = -1
	}

	if bms == nil {
		return uids
	}

	bmsBytes := types.CstrToBytes(bms[:])
	userIDsBytes := bytes.Split(bmsBytes, []byte{'/'})
	userIDs := make([]*ptttype.UserID_t, len(userIDsBytes))
	for idx, each := range userIDsBytes {
		userIDs[idx] = &ptttype.UserID_t{}
		copy(userIDs[idx][:], each)
	}

	// parse user-ids
	idxUID := 0
	for _, each := range userIDs {
		uid, err := SearchUserRaw(each, nil)
		if err != nil || !uid.IsValid() {
			continue
		}
		uids[idxUID] = uid
		idxUID++
	}

	return uids
}

func ResetBoard(bid ptttype.Bid) (err error) {
	if !bid.IsValid() {
		return ptttype.ErrInvalidBid
	}

	bidInCache := bid.ToBidInStore()

	var busystate int32
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BBusyState),
		types.INT32_SZ,
		unsafe.Pointer(&busystate),
	)

	var busystate_b types.Time4
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BusyStateB)+types.TIME4_SZ*uintptr(bidInCache),
		types.TIME4_SZ,
		unsafe.Pointer(&busystate_b),
	)

	nowTS := types.NowTS()
	// busy, return
	if busystate != 0 || nowTS-busystate_b < 10 {
		time.Sleep(1 * time.Second)
		return ErrBusy
	}

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BusyStateB)+types.TIME4_SZ*uintptr(bidInCache),
		types.TIME4_SZ,
		unsafe.Pointer(&nowTS),
	)
	defer func() {
		zeroTS := types.Time4(0)
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.BusyStateB)+types.TIME4_SZ*uintptr(bidInCache),
			types.TIME4_SZ,
			unsafe.Pointer(&zeroTS),
		)
	}()

	file, err := os.Open(ptttype.FN_BOARD)
	if err != nil {
		return err
	}
	defer file.Close()

	board := &ptttype.BoardHeaderRaw{}
	_, err = file.Seek(int64(bidInCache)*int64(ptttype.BOARD_HEADER_RAW_SZ), 0)
	if err != nil {
		return err
	}
	err = types.BinaryRead(file, binary.LittleEndian, board)
	if err != nil {
		return err
	}

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache),
		ptttype.BOARD_HEADER_RAW_SZ,
		unsafe.Pointer(board),
	)

	buildBMCache(bid)

	return nil
}

func buildBMCache(bid ptttype.Bid) {
	if !bid.IsValid() {
		return
	}

	BMs := &ptttype.BM_t{}
	bidInCache := bid.ToBidInStore()
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache)+ptttype.BOARD_HEADER_BM_OFFSET,
		ptttype.BM_SZ,
		unsafe.Pointer(BMs),
	)

	// reset uids
	resetUids := [ptttype.MAX_BMs]ptttype.UID{}
	for idx := 0; idx < ptttype.MAX_BMs; idx++ {
		resetUids[idx] = -1
	}

	const BMCACHE_SZ = unsafe.Sizeof(Shm.Raw.BMCache[0])
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BMCache)+BMCACHE_SZ*uintptr(bidInCache),
		BMCACHE_SZ,
		unsafe.Pointer(&resetUids),
	)

	// set uids
	uids := ParseBMList(BMs)

	if len(uids) == 0 {
		return
	}

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BMCache)+BMCACHE_SZ*uintptr(bidInCache),
		ptttype.UID_SZ*uintptr(ptttype.MAX_BMs),
		unsafe.Pointer(uids),
	)
}

func AddbrdTouchCache() (bid ptttype.Bid, err error) {
	Shm.IncUint32(
		unsafe.Offsetof(Shm.Raw.BNumber),
	)

	nBoards := NumBoards()

	bid = ptttype.Bid(nBoards)

	err = ResetBoard(bid)
	if err != nil {
		return 0, err
	}
	SortBCache()

	return bid, nil
}

func SetLastPosttime(bid ptttype.Bid, nowTS types.Time4) (err error) {
	if !bid.IsValid() {
		return ptttype.ErrInvalidBid
	}

	bidInCache := bid.ToBidInStore()
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.LastPostTime)+types.TIME4_SZ*uintptr(bidInCache),
		types.TIME4_SZ,
		unsafe.Pointer(&nowTS),
	)

	return nil
}

func GetLastPosttime(bid ptttype.Bid) (lastposttime types.Time4, err error) {
	if !bid.IsValid() {
		return 0, ptttype.ErrInvalidBid
	}

	bidInCache := bid.ToBidInStore()
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.LastPostTime)+types.TIME4_SZ*uintptr(bidInCache),
		types.TIME4_SZ,
		unsafe.Pointer(&lastposttime),
	)

	return lastposttime, nil
}

func TouchBPostNum(bid ptttype.Bid, delta int32) (err error) {
	if !bid.IsValid() {
		return ptttype.ErrInvalidBid
	}

	bidInCache := bid.ToBidInStore()
	total := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Total)+types.INT32_SZ*uintptr(bidInCache),
		types.INT32_SZ,
		unsafe.Pointer(&total),
	)

	total += delta

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Total)+types.INT32_SZ*uintptr(bidInCache),
		types.INT32_SZ,
		unsafe.Pointer(&total),
	)

	return nil
}
