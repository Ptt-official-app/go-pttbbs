package ptt

import (
	"syscall"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func isVisibleStat(me *ptttype.UserInfoRaw, uentp *ptttype.UserInfoRaw, friStat ptttype.FriendStat) bool {
	if uentp == nil || uentp.UserID[0] == 0 {
		return false
	}

	if uentp.Mode == ptttype.USER_OP_DEBUGSLEEPING {
		return false
	}

	if uentp.UserLevel.Hide() && !me.UserLevel.Hide() { // 對方紫色隱形而你沒有
		return false
	}

	if me.UserLevel.HasUserPerm(ptttype.PERM_SYSOP) {
		return true
	}

	if friStat.HasPerm(ptttype.FRIEND_STAT_HRM) && friStat.HasPerm(ptttype.FRIEND_STAT_HFM) {
		return true
	}

	if uentp.Invisible && !me.UserLevel.HasUserPerm(ptttype.PERM_SEECLOAK) {
		return false
	}

	return !(friStat.HasPerm(ptttype.FRIEND_STAT_HRM))
}

func friendStat(meID ptttype.UtmpID, me *ptttype.UserInfoRaw, uentID ptttype.UtmpID, uentp *ptttype.UserInfoRaw) (hit ptttype.FriendStat) {
	if me.BrcID != 0 && uentp.BrcID == me.BrcID {
		hit = ptttype.FRIEND_STAT_IBH
	}

	for i := 0; me.FriendOnline[i] != 0 && i < ptttype.MAX_FRIEND; i++ {
		j := me.FriendOnline[i].ToUtmpID()
		if ptttype.ValidUSHMEntry(j) && uentID == j {
			hit |= me.FriendOnline[i].ToFriendStat()
			break
		}
	}

	if uentp.UserLevel.Hide() {
		return hit & ptttype.FRIEND_STAT_FRIEND
	}

	return hit
}

// myWrite
//
// XXX dealing with only WATERBALL_ALOHA for now.
// id in c-pttbbs is the same as puin.UserID
func myWrite(myUtmpID ptttype.UtmpID, myInfo *ptttype.UserInfoRaw, pid types.Pid_t, prompt []byte, flag ptttype.WaterBall, putmpID ptttype.UtmpID, puin *ptttype.UserInfoRaw) (msgCount uint8, err error) {
	if pid >= types.DEFAULT_PID_MAX {
		return 0, nil
	}

	uin := puin
	utmpID := putmpID
	if uin == nil {
		utmpID, uin = cache.SearchUListPID(pid)
		if uin == nil || uin.UserID[0] == 0 {
			return 0, ErrNoUser
		}
	}

	// we don't have water in go-pttbbs.
	if (uin == nil || uin.UserID[0] == 0) && !(flag == ptttype.WATERBALL_GENERAL || flag == ptttype.WATERBALL_ANGEL || flag == ptttype.WATERBALL_ANSWER) {
		return 0, ErrNoUser
	}

	msg := cmsys.StripAnsi(prompt, cmsys.STRIP_ANSI_ALL)

	mode0, c0 := myWriteInit(myUtmpID, myInfo)
	defer myWriteDefer(mode0, c0, myUtmpID, myInfo)

	msgCount, err = myWriteMsg(myUtmpID, myInfo, flag, utmpID, uin, msg)

	return msgCount, err
}

func myWriteInit(myUtmpID ptttype.UtmpID, myInfo *ptttype.UserInfoRaw) (mode ptttype.UserOpMode, c0 byte) {
	mode = cache.Shm.Shm.UInfo[myUtmpID].Mode
	c0 = cache.Shm.Shm.UInfo[myUtmpID].Chatid[0]

	cache.Shm.Shm.UInfo[myUtmpID].Mode = 0
	cache.Shm.Shm.UInfo[myUtmpID].Chatid = ptttype.ChatID_t{}

	return mode, c0
}

func myWriteDefer(mode ptttype.UserOpMode, c0 uint8, myUtmpID ptttype.UtmpID, myInfo *ptttype.UserInfoRaw) {
	cache.Shm.Shm.UInfo[myUtmpID].Mode = mode
	cache.Shm.Shm.UInfo[myUtmpID].Chatid[0] = c0
}

func myWriteMsg(myUtmpID ptttype.UtmpID, myInfo *ptttype.UserInfoRaw, flag ptttype.WaterBall, utmpID ptttype.UtmpID, uin *ptttype.UserInfoRaw, msg []byte) (msgCount uint8, err error) {
	if uin.MsgCount == ptttype.MAX_MSGS-1 {
		return uin.MsgCount, ErrTooManyMsgs
	}

	msgMode := ptttype.MSGMODE_WRITE
	switch flag {
	case ptttype.WATERBALL_ANGEL:
		msgMode = ptttype.MSGMODE_TOANGEL
	case ptttype.WATERBALL_CONFIRM_ANGEL:
		msgMode = ptttype.MSGMODE_TOANGEL
	case ptttype.WATERBALL_ANSWER:
		msgMode = ptttype.MSGMODE_FROMANGEL
	case ptttype.WATERBALL_ALOHA:
		msgMode = ptttype.MSGMODE_ALOHA
	}

	msgQueue := &ptttype.MsgQueueRaw{
		Pid:     myInfo.Pid,
		MsgMode: msgMode,
	}
	copy(msgQueue.UserID[:], myInfo.UserID[:])
	copy(msgQueue.LastCallIn[:], msg)

	msgCount = cache.Shm.Shm.UInfo[utmpID].MsgCount
	if msgCount == ptttype.MAX_MSGS-1 {
		return 0, ErrTooManyMsgs
	}

	idxMsg := msgCount
	msgCount++
	cache.Shm.Shm.UInfo[utmpID].MsgCount++
	cache.Shm.Shm.UInfo[utmpID].Msgs[idxMsg] = *msgQueue

	if ptttype.NOKILLWATERBALL {
		nowTS := types.NowTS()
		cache.Shm.Shm.UInfo[utmpID].WBTime = nowTS
	} else {
		err = types.Kill(uin.Pid, syscall.SIGUSR2)
		if err != nil {
			if flag == ptttype.WATERBALL_ALOHA {
				err = nil
			}
			return msgCount, err
		}
	}

	return msgCount, nil
}
