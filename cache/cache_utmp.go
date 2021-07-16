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
	const offsetUInfo = unsafe.Offsetof(Shm.Raw.UInfo)
	const sizeOfSorted = unsafe.Sizeof(Shm.Raw.Sorted[0])
	const sizeOfSorted2 = unsafe.Sizeof(Shm.Raw.Sorted[0][0])
	userIDInCache := &ptttype.UserID_t{}
	userIDInCache_ptr := unsafe.Pointer(userIDInCache)
	for i := (start + end) / 2; ; i = (start + end) / 2 {
		// get utmpID
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.Sorted)+
				sizeOfSorted*uintptr(currentSorted)+
				sizeOfSorted2*uintptr(ptttype.SORT_BY_ID)+
				ptttype.UTMP_ID_SZ*uintptr(i),
			ptttype.UTMP_ID_SZ,
			utmpID_ptr,
		)

		// get user-id
		Shm.ReadAt(
			offsetUInfo+
				ptttype.USER_INFO_RAW_SZ*uintptr(utmpID)+
				ptttype.USER_INFO_USER_ID_OFFSET,
			ptttype.USER_ID_SZ,
			userIDInCache_ptr,
		)

		// cmp
		j := types.Cstrcasecmp(userID[:], userIDInCache[:])

		if j == 0 {
			uInfo := &ptttype.UserInfoRaw{}
			Shm.ReadAt(
				offsetUInfo+
					ptttype.USER_INFO_RAW_SZ*uintptr(utmpID),
				ptttype.USER_INFO_RAW_SZ,
				unsafe.Pointer(uInfo),
			)

			return utmpID, uInfo
		}

		// determine start / end
		if end == start {
			break
		} else if i == start {
			i = end // nolint
			start = end
		} else if j > 0 {
			start = i
		} else {
			end = i
		}
	}

	return -1, nil
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

	const offsetUInfo = unsafe.Offsetof(Shm.Raw.UInfo)
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

		// get uPid
		Shm.ReadAt(
			offsetUInfo+
				ptttype.USER_INFO_RAW_SZ*uintptr(utmpID)+
				ptttype.USER_INFO_PID_OFFSET,
			types.PID_SZ,
			uPid_ptr,
		)

		// do cmp()
		isDiff = pid - uPid

		if isDiff == 0 {
			uInfo := &ptttype.UserInfoRaw{}
			Shm.ReadAt(
				offsetUInfo+
					ptttype.USER_INFO_RAW_SZ*uintptr(utmpID),
				ptttype.USER_INFO_RAW_SZ,
				unsafe.Pointer(uInfo),
			)

			return utmpID, uInfo
		}

		// determine start / end
		if end == start {
			break
		} else if i == start {
			i = end // nolint
			start = end
		} else if isDiff > 0 {
			start = i
		} else {
			end = i
		}
	}

	return -1, nil
}

//SetUtmpMode
//
//XXX skip utmp for now.
func SetUtmpMode(uid ptttype.UID, mode ptttype.UserOpMode) (err error) {
	/*
		pid := uid.ToPid()
		utmpID, _ := SearchUListPID(pid)
		if utmpID == -1 {
			return ErrNotFound
		}
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.UInfo)+ptttype.USER_INFO_RAW_SZ*uintptr(utmpID)+ptttype.USER_INFO_MODE_OFFSET,
			ptttype.USER_INFO_MODE_SZ,
			unsafe.Pointer(&mode),
		)
	*/

	return nil
}
