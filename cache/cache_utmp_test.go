package cache

import (
	"reflect"
	"sync"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestSearchUListUserID(t *testing.T) {
	setupTest()
	defer teardownTest()

	InitFillUHash(false)

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("not-exists"))

	Shm.Shm.CurrSorted = 0
	Shm.Shm.UTMPNumber = 5

	sorted := []ptttype.UtmpID{2, 1, 0, 4, 3}
	copy(Shm.Shm.Sorted[0][0][:], sorted)

	uinfo := []ptttype.UserInfoRaw{testUserInfo1, testUserInfo2, testUserInfo3, testUserInfo4, testUserInfo5, testUserInfo6}
	copy(Shm.Shm.UInfo[:], uinfo)

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
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			_, got := SearchUListUserID(tt.args.userID)
			testutil.TDeepEqual(t, "userinfo", got, tt.expected)
		})
		wg.Wait()
	}
}

func TestSearchUListPID(t *testing.T) {
	setupTest()
	defer teardownTest()

	InitFillUHash(false)

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("not-exists"))

	Shm.Shm.CurrSorted = 0

	Shm.Shm.UTMPNumber = 5

	sorted := []ptttype.UtmpID{2, 1, 0, 4, 3}
	copy(Shm.Shm.Sorted[0][ptttype.SORT_BY_PID][:], sorted)
	const sizeOfSorted2 = unsafe.Sizeof(Shm.Raw.Sorted[0][0])

	uinfo := []ptttype.UserInfoRaw{testUserInfo1, testUserInfo2, testUserInfo3, testUserInfo4, testUserInfo5, testUserInfo6}
	copy(Shm.Shm.UInfo[:], uinfo)

	type args struct {
		pid types.Pid_t
	}
	tests := []struct {
		name      string
		args      args
		expected  ptttype.UtmpID
		expected1 *ptttype.UserInfoRaw
	}{
		// TODO: Add test cases.
		{
			args:      args{1},
			expected:  2,
			expected1: &testUserInfo3,
		},
		{
			args:      args{2},
			expected:  1,
			expected1: &testUserInfo2,
		},
		{
			args:      args{3},
			expected:  0,
			expected1: &testUserInfo1,
		},
		{
			args:      args{4},
			expected:  4,
			expected1: &testUserInfo5,
		},
		{
			args:      args{5},
			expected:  3,
			expected1: &testUserInfo4,
		},
		{
			args:      args{7},
			expected:  -1,
			expected1: nil,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, got1 := SearchUListPID(tt.args.pid)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SearchUListPID() got = %v, want %v", got, tt.expected)
			}
			if !reflect.DeepEqual(got1, tt.expected1) {
				t.Errorf("SearchUListPID() got1 = %v, want %v", got1, tt.expected1)
			}
		})
		wg.Wait()
	}
}

func TestGetUTotal(t *testing.T) {
	setupTest()
	defer teardownTest()

	InitFillUHash(false)

	Shm.Shm.UTMPNumber = 5
	tests := []struct {
		name      string
		wantTotal int32
	}{
		// TODO: Add test cases.
		{
			"test get user numer",
			5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTotal := GetUTMPNumber(); gotTotal != tt.wantTotal {
				t.Errorf("GetUTotal() = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}
