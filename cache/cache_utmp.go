package cache

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func SearchUListUserID(userID *ptttype.UserID_t) (ptttype.UtmpID, *ptttype.UserInfoRaw) {
	// start and end
	start := int32(0)

	end := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.UTMPNumber),
		types.INT32_SZ,
		unsafe.Pointer(&end),
	)
	end--
	if end < 0 {
		return -1, nil
	}

	// current-sorted (for double-buffer)
	currentSorted := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.CurrSorted),
		types.INT32_SZ,
		unsafe.Pointer(&currentSorted),
	)

	// search
	utmpID := ptttype.UtmpID(0)
	utmpID_p := &utmpID
	utmpID_ptr := unsafe.Pointer(utmpID_p)
	isDiff := 0
	const offsetUInfoUserID = unsafe.Offsetof(Shm.Raw.UInfo[0].UserID)
	const sizeOfSorted = unsafe.Sizeof(Shm.Raw.Sorted[0])
	const sizeOfSorted2 = unsafe.Sizeof(Shm.Raw.Sorted[0][0])
	for i := (start + end) / 2; ; i = (start + end) / 2 {
		// get utmpID
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.Sorted)+sizeOfSorted*uintptr(currentSorted)+sizeOfSorted2*uintptr(ptttype.SORT_BY_ID)+ptttype.UTMP_ID_SZ*uintptr(i),
			ptttype.UTMP_ID_SZ,
			utmpID_ptr,
		)

		// do cmp()
		isDiff = Shm.Memcmp(
			unsafe.Offsetof(Shm.Raw.UInfo)+ptttype.USER_INFO_RAW_SZ*uintptr(utmpID)+offsetUInfoUserID,
			ptttype.USER_ID_SZ,
			unsafe.Pointer(userID),
		)
		//XXX our implementation forces memcmp(shm, userID)
		//    force the same sign as pttbbs (memcmp(userID, shm))
		isDiff = -isDiff

		if isDiff == 0 {
			uInfo := &ptttype.UserInfoRaw{}
			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.UInfo)+ptttype.USER_INFO_RAW_SZ*uintptr(utmpID),
				ptttype.USER_INFO_RAW_SZ,
				unsafe.Pointer(uInfo),
			)

			return utmpID, uInfo
		}

		// determine start / end
		if end == start {
			break
		} else if i == start {
			start = end
		} else if isDiff > 0 {
			start = i
		} else {
			end = i
		}
	}

	return 0, nil
}

func SearchUListPID(pid types.Pid_t) (ptttype.UtmpID, *ptttype.UserInfoRaw) {
	// start and end
	start := int32(0)

	end := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.UTMPNumber),
		types.INT32_SZ,
		unsafe.Pointer(&end),
	)
	end--
	if end < 0 {
		return -1, nil
	}

	// current-sorted (for double-buffer)
	currentSorted := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.CurrSorted),
		types.INT32_SZ,
		unsafe.Pointer(&currentSorted),
	)

	// search
	utmpID := ptttype.UtmpID(0)
	utmpID_p := &utmpID
	utmpID_ptr := unsafe.Pointer(utmpID_p)
	isDiff := types.Pid_t(0)

	uPid := types.Pid_t(0)
	uPid_p := &uPid
	uPid_ptr := unsafe.Pointer(uPid_p)

	const offsetUInfoPid = unsafe.Offsetof(Shm.Raw.UInfo[0].Pid)
	const sizeOfSorted = unsafe.Sizeof(Shm.Raw.Sorted[0])
	const sizeOfSorted2 = unsafe.Sizeof(Shm.Raw.Sorted[0][0])
	for i := (start + end) / 2; ; i = (start + end) / 2 {
		// get utmpID
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.Sorted)+
				sizeOfSorted*uintptr(currentSorted)+
				sizeOfSorted2*uintptr(ptttype.SORT_BY_PID)+
				ptttype.UTMP_ID_SZ*uintptr(i),
			ptttype.UTMP_ID_SZ,
			utmpID_ptr,
		)

		//get uPid
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.UInfo)+
				ptttype.USER_INFO_RAW_SZ*uintptr(utmpID)+
				offsetUInfoPid,
			types.PID_SZ,
			uPid_ptr,
		)

		// do cmp()
		isDiff = pid - uPid

		if isDiff == 0 {
			uInfo := &ptttype.UserInfoRaw{}
			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.UInfo)+ptttype.USER_INFO_RAW_SZ*uintptr(utmpID),
				ptttype.USER_INFO_RAW_SZ,
				unsafe.Pointer(uInfo),
			)

			return utmpID, uInfo
		}

		// determine start / end
		if end == start {
			break
		} else if i == start {
			start = end
		} else if isDiff > 0 {
			start = i
		} else {
			end = i
		}
	}

	return 0, nil
}
