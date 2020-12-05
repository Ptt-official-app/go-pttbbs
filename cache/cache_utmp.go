package cache

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func SearchUListUserID(userID *ptttype.UserID_t) *ptttype.UserInfoRaw {
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
		return nil
	}

	// current-sorted (for double-buffer)
	currentSorted := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.CurrSorted),
		types.INT32_SZ,
		unsafe.Pointer(&currentSorted),
	)

	// search
	uidInCache := ptttype.UidInStore(0)
	pUidInCache := &uidInCache
	uidInCachePtr := unsafe.Pointer(pUidInCache)
	isDiff := 0
	const offsetUInfoUserID = unsafe.Offsetof(Shm.Raw.UInfo[0].UserID)
	const sizeOfSorted = unsafe.Sizeof(Shm.Raw.Sorted[0])
	const sizeOfSorted2 = unsafe.Sizeof(Shm.Raw.Sorted[0][0])
	for i := (start + end) / 2; ; i = (start + end) / 2 {
		// get uidInCache
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.Sorted)+sizeOfSorted*uintptr(currentSorted)+sizeOfSorted2*uintptr(ptttype.SORT_BY_ID)+ptttype.UID_IN_STORE_SZ*uintptr(i),
			types.INT32_SZ,
			uidInCachePtr,
		)

		// do cmp()
		isDiff = Shm.Memcmp(
			unsafe.Offsetof(Shm.Raw.UInfo)+ptttype.USER_INFO_RAW_SZ*uintptr(uidInCache)+offsetUInfoUserID,
			ptttype.USER_ID_SZ,
			unsafe.Pointer(userID),
		)
		//XXX our implementation forces memcmp(shm, userID)
		//    force the same sign as pttbbs (memcmp(userID, shm))
		isDiff = -isDiff

		if isDiff == 0 {
			uInfo := &ptttype.UserInfoRaw{}
			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.UInfo)+ptttype.USER_INFO_RAW_SZ*uintptr(uidInCache),
				ptttype.USER_INFO_RAW_SZ,
				unsafe.Pointer(uInfo),
			)

			return uInfo
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

	return nil
}
