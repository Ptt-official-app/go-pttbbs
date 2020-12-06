package cache

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"time"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func GetBCache(bidInCache ptttype.BidInStore) (board *ptttype.BoardHeaderRaw, err error) {
	board = &ptttype.BoardHeaderRaw{}

	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(bidInCache),
		ptttype.BOARD_HEADER_RAW_SZ,
		unsafe.Pointer(board),
	)
	return board, nil
}

func IsHiddenBoardFriend(bidInCache ptttype.BidInStore, uidInCache ptttype.UidInStore) bool {
	//hbfl time
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

	uid := uidInCache.ToUid()

	var friendID ptttype.Uid
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

	hbfl := [ptttype.MAX_FRIEND + 1]ptttype.Uid{}
	const hbflsz = unsafe.Sizeof(hbfl)

	reader := bufio.NewReader(file)
	var line []byte
	var uid ptttype.Uid
	// num++ is in the end of the for.
	for num := ptttype.Uid(1); num <= ptttype.MAX_FRIEND; {
		line, _ = types.ReadLine(reader)
		if len(line) == 0 {
			break
		}
		theList := bytes.Split(line, []byte{' '}) //The \x00 is taken care of by scanner.

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

		num++ // num++ is in the end of the for.
	}

	hbfl[0] = ptttype.Uid(types.NowTS())

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
	var nboards int32
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.BNumber),
		types.INT32_SZ,
		unsafe.Pointer(&nboards),
	)

	return nboards
}

//Reload BCache
//
//https://github.com/ptt/pttbbs/blob/master/common/bbs/cache.c#L458
func ReloadBCache() {
	var busystate int32
	for i := 0; i < 10; i++ { //Is it ok that we don't use mutex or semaphore here?
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
	// should we check that the busystate is still != 0 and return?

	busystate = 1
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BBusyState),
		types.INT32_SZ,
		unsafe.Pointer(&busystate),
	)

	theBytes, err := reloadBCacheReadFile()
	if err != nil {
		return
	}

	const bcachesz = unsafe.Sizeof(Shm.Raw.BCache)
	var theSize = bcachesz
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

	busystate = 0
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BBusyState),
		types.INT32_SZ,
		unsafe.Pointer(&busystate),
	)

	sortBCache()

	//XXX do not do reloadCacheLoadBottom in test
	if IsTest {
		return
	}

	go reloadCacheLoadBottom()
}

//sortBCache
//XXX TODO: implement
func sortBCache() {
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
	//init vars
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

		var n_uint8 = uint8(n)
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.NBottom)+unsafe.Sizeof(n_uint8)*i,
			1,
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
