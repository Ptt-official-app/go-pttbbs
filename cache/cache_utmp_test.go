package cache

import (
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

func TestSearchUListUserID(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		log.Errorf("TestGetUserID: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	InitFillUHash(false)

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("not-exists"))

	currSorted := 0
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.CurrSorted),
		types.INT32_SZ,
		unsafe.Pointer(&currSorted),
	)

	nUser := 5
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.UTMPNumber),
		types.INT32_SZ,
		unsafe.Pointer(&nUser),
	)

	sorted := [5]ptttype.UidInStore{2, 1, 0, 4, 3}
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Sorted),
		unsafe.Sizeof(sorted),
		unsafe.Pointer(&sorted),
	)

	uinfo := [6]ptttype.UserInfoRaw{testUserInfo1, testUserInfo2, testUserInfo3, testUserInfo4, testUserInfo5, testUserInfo6}
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.UInfo),
		unsafe.Sizeof(uinfo),
		unsafe.Pointer(&uinfo),
	)

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name     string
		args     args
		expected *ptttype.UserInfoRaw
	}{
		// TODO: Add test cases.
		{
			args:     args{userID0},
			expected: nil,
		},
		{
			args:     args{&testUserInfo1.UserID},
			expected: &testUserInfo1,
		},
		{
			args:     args{&testUserInfo2.UserID},
			expected: &testUserInfo2,
		},
		{
			args:     args{&testUserInfo3.UserID},
			expected: &testUserInfo3,
		},
		{
			args:     args{&testUserInfo4.UserID},
			expected: &testUserInfo4,
		},
		{
			args:     args{&testUserInfo5.UserID},
			expected: &testUserInfo5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SearchUListUserID(tt.args.userID)
			types.TDeepEqual(t, got, tt.expected)
		})
	}
}
