package ptt

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func IsBMCache(user *ptttype.UserecRaw, uid ptttype.Uid, bid ptttype.Bid) bool {
	bidInCache := bid.ToBidInStore()

	// XXX potential issue: (thanks for mtdas@ptt)
	//  buildBMcache use -1 as "none".
	//  some function may call is_BM_cache early
	//  without having currutmp->uid (maybe?)
	//  and may get BM permission accidentally.
	// quick check

	if !user.UserLevel.HasUserPerm(ptttype.PERM_BASIC) ||
		uid == 0 ||
		uid == -1 {
		return false
	}
	if !user.UserLevel.HasBasicUserPerm(ptttype.PERM_LOGINCLOAK) {
		return false
	}

	pbm := [ptttype.MAX_BMs]ptttype.Uid{}
	const bmcache0sz = unsafe.Sizeof(cache.Shm.Raw.BMCache[0])
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BMCache)+uintptr(bidInCache)*bmcache0sz,
		bmcache0sz,
		unsafe.Pointer(&pbm),
	)
	if uid == pbm[0] || uid == pbm[1] || uid == pbm[2] || uid == pbm[3] {
		if user.UserLevel.HasUserPerm(ptttype.PERM_BM) {
			pwcuBitEnableLevel(uid, &user.UserID, ptttype.PERM_BM)
		}
		return true
	}

	return false
}
