package ptt

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func Test_isVisibleStat(t *testing.T) {
	setupTest()
	defer teardownTest()

	userID := ptttype.UserID_t{}
	copy(userID[:], []byte("SYSOP"))

	uent1 := &ptttype.UserInfoRaw{
		UserID: userID,
		Mode:   ptttype.USER_OP_DEBUGSLEEPING,
	}

	me2 := &ptttype.UserInfoRaw{
		UserLevel: ptttype.PERM_SYSOP,
	}

	uent3 := &ptttype.UserInfoRaw{
		UserID:    userID,
		UserLevel: ptttype.PERM_SYSOPHIDE,
	}

	me3 := &ptttype.UserInfoRaw{
		UserLevel: ptttype.PERM_BASIC,
	}

	uent4 := &ptttype.UserInfoRaw{
		UserID:    userID,
		Invisible: true,
	}

	uent5 := &ptttype.UserInfoRaw{
		UserID: userID,
	}

	type args struct {
		me      *ptttype.UserInfoRaw
		uentp   *ptttype.UserInfoRaw
		friStat ptttype.FriendStat
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			expected: false,
		},
		{
			args:     args{me: me2, uentp: uent1},
			expected: false,
		},
		{
			args:     args{me: me2, uentp: uent3},
			expected: false,
		},
		{
			name:     "check with PERM_SYSOP",
			args:     args{me: me2, uentp: uent4},
			expected: true,
		},
		{
			args:     args{me: me3, uentp: uent3},
			expected: false,
		},
		{
			args:     args{friStat: ptttype.FRIEND_STAT_HRM | ptttype.FRIEND_STAT_HFM},
			expected: false,
		},
		{
			args:     args{me: me3, uentp: uent4},
			expected: false,
		},
		{
			args:     args{me: me3, uentp: uent5},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isVisibleStat(tt.args.me, tt.args.uentp, tt.args.friStat); got != tt.expected {
				t.Errorf("isVisibleStat() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_friendStat(t *testing.T) {
	setupTest()
	defer teardownTest()

	me0 := &ptttype.UserInfoRaw{
		BrcID: 1,
	}
	uent0 := &ptttype.UserInfoRaw{
		BrcID: 1,
	}

	me1 := &ptttype.UserInfoRaw{}
	theFriendStat := ptttype.FRIEND_STAT_HFM | ptttype.FRIEND_STAT_IFH
	me1.FriendOnline[0] = ptttype.FriendOnline(12 + (theFriendStat << 24))
	uent1 := &ptttype.UserInfoRaw{}

	type args struct {
		meID   ptttype.UtmpID
		me     *ptttype.UserInfoRaw
		uentID ptttype.UtmpID
		uentp  *ptttype.UserInfoRaw
	}
	tests := []struct {
		name        string
		args        args
		expectedHit ptttype.FriendStat
	}{
		// TODO: Add test cases.
		{
			args:        args{me: me0, uentp: uent0},
			expectedHit: ptttype.FRIEND_STAT_IBH,
		},
		{
			args:        args{me: me1, uentp: uent1, uentID: 12},
			expectedHit: ptttype.FRIEND_STAT_HFM | ptttype.FRIEND_STAT_IFH,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHit := friendStat(tt.args.meID, tt.args.me, tt.args.uentID, tt.args.uentp); !reflect.DeepEqual(gotHit, tt.expectedHit) {
				t.Errorf("friendStat() = %v, want %v", gotHit, tt.expectedHit)
			}
		})
	}
}

func Test_myWrite(t *testing.T) {
	setupTest()
	defer teardownTest()

	currSorted := 0
	cache.Shm.WriteAt(
		unsafe.Offsetof(cache.Shm.Raw.CurrSorted),
		types.INT32_SZ,
		unsafe.Pointer(&currSorted),
	)

	nUser := 5
	cache.Shm.WriteAt(
		unsafe.Offsetof(cache.Shm.Raw.UTMPNumber),
		types.INT32_SZ,
		unsafe.Pointer(&nUser),
	)

	sorted := [6]ptttype.UtmpID{2, 1, 0, 4, 3}
	const sizeOfSorted2 = unsafe.Sizeof(cache.Shm.Raw.Sorted[0][0])
	cache.Shm.WriteAt(
		unsafe.Offsetof(cache.Shm.Raw.Sorted)+sizeOfSorted2*uintptr(ptttype.SORT_BY_PID),
		unsafe.Sizeof(sorted),
		unsafe.Pointer(&sorted),
	)

	uinfo := [6]ptttype.UserInfoRaw{testUserInfo1, testUserInfo2, testUserInfo3, testUserInfo4, testUserInfo5, testUserInfo6}
	cache.Shm.WriteAt(
		unsafe.Offsetof(cache.Shm.Raw.UInfo),
		unsafe.Sizeof(uinfo),
		unsafe.Pointer(&uinfo),
	)

	prompt0 := []byte("test")

	type args struct {
		myUtmpID ptttype.UtmpID
		myInfo   *ptttype.UserInfoRaw
		pid      types.Pid_t
		prompt   []byte
		flag     ptttype.WaterBall
		putmpID  ptttype.UtmpID
		puin     *ptttype.UserInfoRaw
	}
	tests := []struct {
		name             string
		args             args
		expectedMsgCount uint8
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			args:             args{myUtmpID: 0, myInfo: &testUserInfo1, pid: 5, prompt: prompt0, flag: ptttype.WATERBALL_ALOHA},
			expectedMsgCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsgCount, err := myWrite(tt.args.myUtmpID, tt.args.myInfo, tt.args.pid, tt.args.prompt, tt.args.flag, tt.args.putmpID, tt.args.puin)
			if (err != nil) != tt.wantErr {
				t.Errorf("myWrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMsgCount != tt.expectedMsgCount {
				t.Errorf("myWrite() = %v, want %v", gotMsgCount, tt.expectedMsgCount)
			}
		})
	}
}
