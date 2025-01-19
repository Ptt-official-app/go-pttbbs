package cache

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"strings"
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

	return &SHM.Shm.BCache[bidInCache], nil
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

	return int32(SHM.Shm.NBottom[bidInCache])
}

func GetBottomTotalAllGuest(boardID *ptttype.BoardID_t) (total int32) {
	ret, ok := MAP.BoardNBottom[*boardID]
	nowUS := types.NowUS()

	if !ok || (ret.count == 0 && nowUS-ret.updateTimeUS > BOARD_ZERO_BOTTOM_RECHECK_TIME_US) || nowUS-ret.updateTimeUS > BOARD_BOTTOM_RECHECK_TIME_US {
		err := SetBottomTotalAllGuest(boardID)
		if err != nil {
			return 0
		}

		ret, ok = MAP.BoardNBottom[*boardID]
		if !ok {
			return 0
		}
	}

	return ret.count
}

func GetBTotal(bid ptttype.Bid) (total int32) {
	if !bid.IsValid() {
		return 0
	}

	bidInCache := bid.ToBidInStore()

	return SHM.Shm.Total[bidInCache]
}

func GetBTotalAllGuest(boardID *ptttype.BoardID_t) (total int32) {
	ret, ok := MAP.BoardTotal[*boardID]
	nowUS := types.NowUS()

	if !ok || (ret.count == 0 && nowUS-ret.updateTimeUS > BOARD_ZERO_TOTAL_RECHECK_TIME_US) || nowUS-ret.updateTimeUS > BOARD_TOTAL_RECHECK_TIME_US {
		err := SetBTotalAllGuest(boardID)
		if err != nil {
			return 0
		}

		ret, ok = MAP.BoardTotal[*boardID]
		if !ok {
			return 0
		}
	}

	return ret.count
}

// SetBTotal
//
// It's possible that we loaded nothing from ReloadBCache in the beginning of the program, and then there are some articles after a while.
// We need to sync the btotal and lastposttime back to shm.
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
	SHM.Shm.Total[bidInCache] = nArticles

	if nArticles == 0 {
		SHM.Shm.LastPostTime[bidInCache] = 0
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
		SHM.Shm.LastPostTime[bidInCache] = 0
		return nil
	}

	createTime, err := articleFilename.CreateTime()
	if err != nil {
		return err
	}

	SHM.Shm.LastPostTime[bidInCache] = createTime

	return nil
}

// SetBTotal
//
// It's possible that we loaded nothing from ReloadBCache in the beginning of the program, and then there are some articles after a while.
// We need to sync the btotal and lastposttime back to shm.
func SetBTotalAllGuest(boardID *ptttype.BoardID_t) (err error) {
	dirFilename, err := path.SetBFile(boardID, ptttype.FN_DIR)
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
	nowUS := types.NowUS()
	MAP.BoardTotal[*boardID] = mapCount{count: nArticles, updateTimeUS: nowUS}

	if nArticles == 0 {
		MAP.BoardLastPostTime[*boardID] = mapTime{time: 0, updateTimeUS: nowUS}
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
		MAP.BoardLastPostTime[*boardID] = mapTime{time: 0, updateTimeUS: nowUS}
		return nil
	}

	createTime, err := articleFilename.CreateTime()
	if err != nil {
		return err
	}

	MAP.BoardLastPostTime[*boardID] = mapTime{time: createTime, updateTimeUS: nowUS}

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
	n := uint8(cmsys.GetNumRecords(bottomFilename, ptttype.FILE_HEADER_RAW_SZ))
	if n > 5 {
		_ = syscall.Unlink(bottomFilename)
		SHM.Shm.NBottom[bidInCache] = 0

		return nil
	}

	SHM.Shm.NBottom[bidInCache] = n

	return nil
}

func SetBottomTotalAllGuest(boardID *ptttype.BoardID_t) error {
	bottomFilename, err := path.SetBFile(boardID, ptttype.FN_DIR_BOTTOM)
	if err != nil {
		return err
	}

	n := int32(cmsys.GetNumRecords(bottomFilename, ptttype.FILE_HEADER_RAW_SZ))
	nowUS := types.NowUS()
	if n > 5 {
		_ = syscall.Unlink(bottomFilename)
		MAP.BoardNBottom[*boardID] = mapCount{count: 0, updateTimeUS: nowUS}

		return nil
	}

	MAP.BoardNBottom[*boardID] = mapCount{count: n, updateTimeUS: nowUS}

	return nil
}

func IsHiddenBoardFriend(bidInCache ptttype.BidInStore, uidInCache ptttype.UIDInStore) bool {
	if types.IS_ALL_GUEST {
		return false
	}

	if !bidInCache.ToBid().IsValid() || !uidInCache.ToUID().IsValid() {
		return false
	}

	// hbfl time
	loadTime := types.Time4(SHM.Shm.Hbfl[bidInCache][0])

	// XXX use nowTS to replace loginStartTime.
	//     HBFLexpire is set as 5-days. nowTS should be ok.
	nowTS := types.NowTS()
	if loadTime < nowTS-types.Time4(ptttype.HBFLexpire) {
		HbflReload(bidInCache)
	}

	uid := uidInCache.ToUID()

	var friendID ptttype.UID
	for i := uintptr(1); i <= ptttype.MAX_FRIEND; i++ {
		friendID = SHM.Shm.Hbfl[bidInCache][i]
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
	if types.IS_ALL_GUEST {
		return
	}

	if !bidInCache.ToBid().IsValid() {
		return
	}

	brdname := &SHM.Shm.BCache[bidInCache].Brdname
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

	copy(SHM.Shm.Hbfl[bidInCache][:], hbfl[:])
}

// NumBoards
//
// https://github.com/ptt/pttbbs/blob/master/common/bbs/cache.c#L512
func NumBoards() int32 {
	return SHM.GetBNumber()
}

func NHots() (nhots uint8) {
	if types.IS_ALL_GUEST {
		return 0
	}

	return SHM.Shm.NHOTs
}

// Reload BCache
//
// https://github.com/ptt/pttbbs/blob/master/common/bbs/cache.c#L458
func ReloadBCache() {
	if types.IS_ALL_GUEST {
		return
	}

	for i := 0; i < 10; i++ { // Is it ok that we don't use mutex or semaphore here?
		if SHM.Shm.BBusyState == 0 {
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
	SHM.Shm.BBusyState = 1
	defer func() {
		SHM.Shm.BBusyState = 0
	}()

	theBytes, err := reloadBCacheReadFile()
	if err != nil {
		return
	}

	const bcachesz = unsafe.Sizeof(SHM.Raw.BCache)
	bcacheBytes := (*[bcachesz]byte)(unsafe.Pointer(&SHM.Shm.BCache))
	copy(bcacheBytes[:], theBytes)

	theSize := uintptr(len(theBytes))
	if bcachesz < theSize {
		theSize = bcachesz
	}
	SHM.Shm.BNumber = int32(theSize / ptttype.BOARD_HEADER_RAW_SZ)

	for i := 0; i < ptttype.MAX_BOARD; i++ {
		SHM.Shm.LastPostTime[i] = 0
		SHM.Shm.Total[i] = 0
	}
	SHM.Shm.BTouchTime = SHM.Shm.BUptime
}

func ReloadBCacheAllGuest() {
	rootPath := strings.Join([]string{
		ptttype.BBSHOME,
		ptttype.DIR_BOARD,
	}, string(os.PathSeparator))

	resetBoardID := &ptttype.BoardID_t{}
	boardID := &ptttype.BoardID_t{}

	MAP.ReloadStartTimeUS = types.NowUS()
	defer func() {
		MAP.ReloadEndTimeUS = types.NowUS()
	}()

	for _, dir1 := range ptttype.BOARD_PREFIXES {
		fullDir1 := strings.Join([]string{rootPath, dir1}, string(os.PathSeparator))
		dirEntries, err := os.ReadDir(fullDir1)
		if err != nil {
			continue
		}

		for _, dirEntry := range dirEntries {
			if !dirEntry.IsDir() {
				continue
			}

			dir2 := dirEntry.Name()
			if !strings.HasPrefix(dir2, dir1) {
				continue
			}

			copy(boardID[:], resetBoardID[:])
			copy(boardID[:], dir2[:])

			_ = SetBTotalAllGuest(boardID)

			_ = SetBottomTotalAllGuest(boardID)
		}
	}
}

func CronReloadBCacheAllGuest() {
	for {
		ReloadBCacheAllGuest()
		time.Sleep(CRON_RELOAD_BCACHE_DURATION)
	}
}

// SortBCache
// XXX TODO: implement
func SortBCache() {
	if types.IS_ALL_GUEST {
		return
	}

	if SHM.Shm.BBusyState != 0 {
		time.Sleep(1 * time.Second)
		return
	}

	SHM.Shm.BBusyState = 1
	defer func() {
		SHM.Shm.BBusyState = 0
	}()

	// init is in shm.c
	SHM.QsortCmpBoardName()
	SHM.QsortCmpBoardClass()

	// for-loop cleaning first-child
	// init vars
	bnumber := SHM.GetBNumber()
	for i := int32(0); i < bnumber; i++ {
		for j := 0; j < int(ptttype.BSORT_BY_MAX); j++ {
			SHM.Shm.BCache[i].FirstChild[j] = 0
		}
	}
}

func reloadCacheLoadBottom() {
	var boardName *ptttype.BoardID_t
	for i := uintptr(0); i < ptttype.MAX_BOARD; i++ {
		boardName = &SHM.Shm.BCache[i].Brdname
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

		SHM.Shm.NBottom[i] = uint8(n)
	}
}

func reloadBCacheReadFile() ([]byte, error) {
	file, err := os.Open(ptttype.FN_BOARD)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	theBytes, err := io.ReadAll(file)
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
	if SHM.Shm.BBusyState != 0 {
		time.Sleep(1 * time.Second)
	}

	// start and end
	start := int32(0)
	end := SHM.GetBNumber()
	end--
	if end < 0 { // unable to get bid
		return -1, 0, nil
	}

	bidInCache := ptttype.BidInStore(0)
	var boardIDInCache *ptttype.BoardID_t
	idx_i32 := (start + end) / 2
	for ; ; idx_i32 = (start + end) / 2 {
		bidInCache = SHM.Shm.BSorted[ptttype.BSORT_BY_NAME][idx_i32]
		boardIDInCache = &SHM.Shm.BCache[bidInCache].Brdname

		j := types.Cstrcasecmp(boardID[:], boardIDInCache[:])
		if j == 0 {
			bid = bidInCache.ToBid()
			return ptttype.SortIdxInStore(idx_i32), bid, nil
		}

		if end == start {
			break
		} else if idx_i32 == start {
			idx_i32 = end //nolint
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
	if SHM.Shm.BBusyState != 0 {
		time.Sleep(1 * time.Second)
	}

	// start and end
	start := int32(0)
	end := SHM.GetBNumber()
	end--
	if end < 0 { // unable to get bid
		return -1, 0, nil
	}

	bidInCache := ptttype.BidInStore(0)

	var titleInCache *ptttype.BoardTitle_t
	var clsInCache []byte
	var boardIDInCache *ptttype.BoardID_t
	idx_i32 := (start + end) / 2
	for ; ; idx_i32 = (start + end) / 2 {
		bidInCache = SHM.Shm.BSorted[ptttype.BSORT_BY_CLASS][idx_i32]
		titleInCache = &SHM.Shm.BCache[bidInCache].Title
		clsInCache = titleInCache.BoardClass()
		boardIDInCache = &SHM.Shm.BCache[bidInCache].Brdname

		j := cmpBoardByClass(cls, boardID, clsInCache, boardIDInCache)
		if j == 0 {
			bid = bidInCache.ToBid()
			return ptttype.SortIdxInStore(idx_i32), bid, nil
		}

		if end == start {
			break
		} else if idx_i32 == start {
			idx_i32 = end //nolint
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

	nBoard_i32 := SHM.GetBNumber()
	nBoard := ptttype.SortIdxInStore(nBoard_i32)
	bidInCache := ptttype.BidInStore(0)
	boardIDInCache := (*ptttype.BoardID_t)(nil)
	if isAsc {
		for ; idxInStore < nBoard; idxInStore++ {
			bidInCache = SHM.Shm.BSorted[ptttype.BSORT_BY_NAME][idxInStore]
			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}
			boardIDInCache = &SHM.Shm.BCache[bidInCache].Brdname

			j := types.Cstrcasecmp(boardID[:], boardIDInCache[:])
			if j <= 0 {
				break
			}
		}
		if idxInStore == nBoard {
			idxInStore = -1
		}
	} else {
		for ; idxInStore >= 0; idxInStore-- {
			bidInCache = SHM.Shm.BSorted[ptttype.BSORT_BY_NAME][idxInStore]
			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}
			boardIDInCache = &SHM.Shm.BCache[bidInCache].Brdname

			j := types.Cstrcasecmp(boardID[:], boardIDInCache[:])
			if j >= 0 {
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

	nBoard_i32 := SHM.GetBNumber()
	nBoard := ptttype.SortIdxInStore(nBoard_i32)
	bidInCache := ptttype.BidInStore(0)
	titleInCache := (*ptttype.BoardTitle_t)(nil)
	boardIDInCache := (*ptttype.BoardID_t)(nil)
	if isAsc {
		for ; idxInStore < nBoard; idxInStore++ {
			bidInCache = SHM.Shm.BSorted[ptttype.BSORT_BY_CLASS][idxInStore]
			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}
			titleInCache = &SHM.Shm.BCache[bidInCache].Title
			clsInCache := titleInCache.BoardClass()
			boardIDInCache = &SHM.Shm.BCache[bidInCache].Brdname

			j := cmpBoardByClass(cls, boardID, clsInCache, boardIDInCache)
			if j <= 0 {
				break
			}
		}
		if idxInStore == nBoard {
			idxInStore = -1
		}
	} else {
		for ; idxInStore >= 0; idxInStore-- {
			bidInCache = SHM.Shm.BSorted[ptttype.BSORT_BY_CLASS][idxInStore]
			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}
			titleInCache = &SHM.Shm.BCache[bidInCache].Title
			clsInCache := titleInCache.BoardClass()
			boardIDInCache = &SHM.Shm.BCache[bidInCache].Brdname

			j := cmpBoardByClass(cls, boardID, clsInCache, boardIDInCache)
			if j >= 0 {
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
	if types.IS_ALL_GUEST {
		return 0, nil
	}

	boardID := findBoardClosetKeyword(keyword, isAsc)
	nBoard_i32 := SHM.GetBNumber()
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
	bidInCache := ptttype.BidInStore(0)
	boardIDInCache := (*ptttype.BoardID_t)(nil)
	// it should be either current idx or the next idx
	const MAX_ITER_FIND_AUTO_COMPLETE = 3
	if isAsc {
		i := 0
		for ; i < MAX_ITER_FIND_AUTO_COMPLETE && idxInStore < nBoard; i, idxInStore = i+1, idxInStore+1 {
			bidInCache = SHM.Shm.BSorted[ptttype.BSORT_BY_NAME][idxInStore]
			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}
			boardIDInCache = &SHM.Shm.BCache[bidInCache].Brdname
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
			bidInCache = SHM.Shm.BSorted[ptttype.BSORT_BY_NAME][idxInStore]
			if !bidInCache.ToBid().IsValid() {
				idxInStore = -1
				break
			}
			boardIDInCache = &SHM.Shm.BCache[bidInCache].Brdname
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
	nowTS := types.NowTS()
	// busy, return
	if SHM.Shm.BBusyState != 0 || nowTS-SHM.Shm.BusyStateB[bidInCache] < 10 {
		time.Sleep(1 * time.Second)
		return ErrBusy
	}

	SHM.Shm.BusyStateB[bidInCache] = nowTS
	defer func() {
		SHM.Shm.BusyStateB[bidInCache] = 0
	}()

	file, err := os.Open(ptttype.FN_BOARD)
	if err != nil {
		return err
	}
	defer file.Close()

	board := &SHM.Shm.BCache[bidInCache]
	_, err = file.Seek(int64(bidInCache)*int64(ptttype.BOARD_HEADER_RAW_SZ), 0)
	if err != nil {
		return err
	}
	err = types.BinaryRead(file, binary.LittleEndian, board)
	if err != nil {
		return err
	}

	buildBMCache(bid)

	return nil
}

func buildBMCache(bid ptttype.Bid) {
	if !bid.IsValid() {
		return
	}

	bidInCache := bid.ToBidInStore()
	BMs := &SHM.Shm.BCache[bidInCache].BM

	// reset uids
	resetUids := [ptttype.MAX_BMs]ptttype.UID{}
	for idx := 0; idx < ptttype.MAX_BMs; idx++ {
		resetUids[idx] = -1
	}

	SHM.Shm.BMCache[bidInCache] = resetUids

	// set uids
	uids := ParseBMList(BMs)
	SHM.Shm.BMCache[bidInCache] = *uids
}

func AddbrdTouchCache() (bid ptttype.Bid, err error) {
	SHM.Shm.BNumber++

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
	SHM.Shm.LastPostTime[bidInCache] = nowTS

	return nil
}

func GetLastPosttime(bid ptttype.Bid) (lastposttime types.Time4, err error) {
	if !bid.IsValid() {
		return 0, ptttype.ErrInvalidBid
	}

	bidInCache := bid.ToBidInStore()

	return SHM.Shm.LastPostTime[bidInCache], nil
}

func TouchBPostNum(bid ptttype.Bid, delta int32) (err error) {
	if !bid.IsValid() {
		return ptttype.ErrInvalidBid
	}

	bidInCache := bid.ToBidInStore()
	SHM.Shm.Total[bidInCache] += delta

	return nil
}

func ResolveBoardGroup(gid ptttype.Bid, bsortBy ptttype.BSortBy) (err error) {
	boardCount := int(NumBoards())
	if boardCount < 0 || boardCount > ptttype.MAX_BOARD {
		return ErrInvalidNumBoards
	}

	bidInCache := ptttype.BidInStore(0)
	parentBoard, err := GetBCache(gid)
	if err != nil {
		return err
	}
	currentBoard := parentBoard
	currentBid := gid
	for idxInStore := 0; idxInStore < boardCount; idxInStore++ {
		bidInCache = SHM.Shm.BSorted[bsortBy][idxInStore]
		bid := bidInCache.ToBid()
		if !bid.IsValid() {
			continue
		}
		board, err := GetBCache(bid)
		if err != nil {
			continue
		}
		if board.Brdname[0] == '\x00' {
			continue
		}
		if board.Gid != gid {
			continue
		}
		if currentBoard == parentBoard {
			err = setBoardFirstChild(currentBid, bsortBy, bid)
			if err != nil {
				return err
			}
		} else {
			err = setBoardNextChild(currentBid, bsortBy, bid)
			if err != nil {
				return err
			}
			err = setBoardParent(currentBid, gid)
			if err != nil {
				return err
			}
		}

		currentBoard = board
		currentBid = bid
	}

	return nil
}

func setBoardFirstChild(bid ptttype.Bid, bsortBy ptttype.BSortBy, childBid ptttype.Bid) (err error) {
	bidInCache := bid.ToBidInStore()
	SHM.Shm.BCache[bidInCache].FirstChild[bsortBy] = childBid

	return nil
}

func setBoardNextChild(bid ptttype.Bid, bsortBy ptttype.BSortBy, childBid ptttype.Bid) (err error) {
	bidInCache := bid.ToBidInStore()
	SHM.Shm.BCache[bidInCache].Next[bsortBy] = childBid

	return nil
}

func setBoardParent(bid ptttype.Bid, parentBid ptttype.Bid) (err error) {
	bidInCache := bid.ToBidInStore()
	SHM.Shm.BCache[bidInCache].Parent = parentBid

	return nil
}

func SetBoardChildCount(bid ptttype.Bid, count int32) (err error) {
	bidInCache := bid.ToBidInStore()
	SHM.Shm.BCache[bidInCache].ChildCount = count

	return nil
}
