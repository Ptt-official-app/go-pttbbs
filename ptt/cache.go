package ptt

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
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
	if !user.UserLevel.HasBasicUserPerm(ptttype.PERM_LOGINOK) {
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

func GetUserLevel(userID *ptttype.UserID_t) (userLevel ptttype.PERM, err error) {
	uid, err := cache.SearchUserRaw(userID, nil)
	if err != nil {
		return 0, err
	}
	if !uid.IsValid() {
		return 0, ptttype.ErrInvalidUserID
	}

	return cmbbs.PasswdQueryUserLevel(uid)
}

func GetUser2(userID *ptttype.UserID_t) (user *ptttype.Userec2Raw, err error) {
	uid, err := cache.SearchUserRaw(userID, nil)
	if err != nil {
		return nil, err
	}
	if !uid.IsValid() {
		return nil, ptttype.ErrInvalidUserID
	}

	user, err = cmbbs.PasswdGetUser2(userID)
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

//postpermMsg
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/cache.c#L209
func postpermMsg(uid ptttype.Uid, user *ptttype.UserecRaw, bid ptttype.Bid, board *ptttype.BoardHeaderRaw) (err error) {
	if isReadonlyBoard(&board.Brdname) {
		return ErrReadOnly
	}

	if user.UserLevel.HasUserPerm(ptttype.PERM_SYSOP) {
		return nil
	}

	err = bannedMsg(user, board)
	if err != nil {
		return err
	}

	if types.Cstrcmp(board.Brdname[:], ptttype.DEFAULT_BOARD) == 0 {
		return nil
	}

	if board.BrdAttr.HasPerm(ptttype.BRD_GUESTPOST) {
		return nil
	}

	if !user.UserLevel.HasUserPerm(ptttype.PERM_POST) {
		return ErrPermitNoPost
	}

	//秘密看板特別處理.
	if board.BrdAttr.HasPerm(ptttype.BRD_HIDE) {
		return nil
	}

	if board.BrdAttr.HasPerm(ptttype.BRD_RESTRICTEDPOST) && !cache.IsHiddenBoardFriend(bid.ToBidInStore(), uid.ToUidInStore()) {
		return ErrRestricted
	}

	if user.UserLevel.HasUserPerm(ptttype.PERM_VIOLATELAW) {
		if board.Level.HasUserPerm(ptttype.PERM_VIOLATELAW) {
			return nil
		} else {
			return ErrViolateLaw
		}
	}

	//除了"post"以外的其他權限要求
	requiredLevel := board.Level & ^ptttype.PERM_POST
	if requiredLevel == 0 {
		return nil
	}
	if !user.UserLevel.HasUserPerm(requiredLevel) {
		return ErrNotPermitted
	}

	return nil
}

func bannedMsg(user *ptttype.UserecRaw, board *ptttype.BoardHeaderRaw) (err error) {
	if ptttype.USE_NEW_BAN_SYSTEM {
		expireTS, _ := isBannedByBoard(user, board)
		nowTS := types.NowTS()
		if expireTS > nowTS {
			return ErrBanned
		}
	} else {
		filename, err := path.SetBFile(&board.Brdname, ptttype.FN_WATER)
		if err != nil {
			return err
		}

		if cmsys.FileExistsRecord(filename, user.UserID[:]) {
			return ErrBanned
		}

	}
	return nil
}
