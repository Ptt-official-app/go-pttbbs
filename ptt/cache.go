package ptt

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
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

func GetUser(userID *ptttype.UserID_t) (user *ptttype.UserecRaw, err error) {
	uid, err := cache.SearchUserRaw(userID, nil)
	if err != nil {
		return nil, err
	}
	if !uid.IsValid() {
		return nil, ptttype.ErrInvalidUserID
	}

	//passwdSyncQuery includes cache.MoneyOf
	user, err = passwdSyncQuery(uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getNewUtmpEnt(uinfo *ptttype.UserInfoRaw) (utmpID ptttype.UtmpID, err error) {
	p := cmsys.StringHash(uinfo.UserID[:]) % ptttype.USHM_SIZE

	var pid types.Pid_t
	ppid := &pid

	for idx := 0; idx < ptttype.USHM_SIZE; idx, p = idx+1, p+1 {
		cache.Shm.ReadAt(
			unsafe.Offsetof(cache.Shm.Raw.UInfo)+uintptr(p)*ptttype.USER_INFO_RAW_SZ+unsafe.Offsetof(ptttype.EMPTY_USER_INFO_RAW.Pid),
			types.PID_SZ,
			unsafe.Pointer(ppid),
		)
		//found same pid.
		//update the newest status.
		//XXX race condition with auto-logout.
		//XXX c-pttbbs does not care the race-condition here.
		//XXX we may not do anything with utmpID though.
		if pid == uinfo.Pid {
			cache.Shm.WriteAt(
				unsafe.Offsetof(cache.Shm.Raw.UInfo)+uintptr(p)*ptttype.USER_INFO_RAW_SZ,
				ptttype.USER_INFO_RAW_SZ,
				unsafe.Pointer(uinfo),
			)

			//https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L998
			one := uint8(1)
			cache.Shm.WriteAt(
				unsafe.Offsetof(cache.Shm.Raw.UTMPNeedSort),
				types.UINT8_SZ,
				unsafe.Pointer(&one),
			)

			return ptttype.UtmpID(p), nil
		}

		//new pid
		if pid == 0 {
			cache.Shm.WriteAt(
				unsafe.Offsetof(cache.Shm.Raw.UInfo)+uintptr(p)*ptttype.USER_INFO_RAW_SZ,
				ptttype.USER_INFO_RAW_SZ,
				unsafe.Pointer(uinfo),
			)

			//https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L998
			one := uint8(1)
			cache.Shm.WriteAt(
				unsafe.Offsetof(cache.Shm.Raw.UTMPNeedSort),
				types.UINT8_SZ,
				unsafe.Pointer(&one),
			)

			return ptttype.UtmpID(p), nil
		}
	}

	return ptttype.UtmpID(-1), ErrNewUtmp
}
