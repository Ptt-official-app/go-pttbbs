package ptt

import (
	"reflect"
	"sync"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func TestLoginQuery(t *testing.T) {
	setupTest()
	defer teardownTest()

	userid1 := ptttype.UserID_t{}
	copy(userid1[:], []byte("SYSOP"))

	type args struct {
		userID *ptttype.UserID_t
		passwd []byte
		ip     *ptttype.IPv4_t
	}
	tests := []struct {
		name        string
		args        args
		expected    *ptttype.UserecRaw
		expectedUid ptttype.Uid
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{userID: &userid1, passwd: []byte("123123")},
			expected:    testUserecRaw1,
			expectedUid: 1,
		},
		{
			args:     args{userID: &userid1, passwd: []byte("124")},
			expected: nil,
			wantErr:  true,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			gotUid, got, err := LoginQuery(tt.args.userID, tt.args.passwd, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUid != tt.expectedUid {
				t.Errorf("LoginQuery() uid: %v expected: %v", gotUid, tt.expectedUid)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("LoginQuery() = %v, expected %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func Test_newUserInfoRaw(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uid  ptttype.Uid
		user *ptttype.UserecRaw
		ip   *ptttype.IPv4_t
		op   ptttype.UserOpMode
	}
	tests := []struct {
		name                 string
		args                 args
		expected             *ptttype.UserInfoRaw
		expectedOrigNickname ptttype.Nickname_t
	}{
		// TODO: Add test cases.
		{
			args: args{
				uid:  10,
				user: testNewUserInfoRawUserecRaw,
				ip:   &ptttype.IPv4_t{'1', '9', '2', '.', '1', '6', '8', '.', '0', '.', '1'},
				op:   ptttype.USER_OP_LOGIN,
			},
			expected:             testNewUserInfoRawUserInfoRaw,
			expectedOrigNickname: testNewUserInfoRawNickname,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			logrus.Infof("tt.args.ip: %v", tt.args.ip)
			got := newUserInfoRaw(tt.args.uid, tt.args.user, tt.args.ip, tt.args.op)
			tt.expected.LastAct = got.LastAct
			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
		wg.Wait()
	}
}

func TestLogin(t *testing.T) {
	setupTest()
	defer teardownTest()

	SetupNewUser(testSetupNewUser1)

	userid1 := ptttype.UserID_t{}
	copy(userid1[:], []byte("A0"))

	ip1 := &ptttype.IPv4_t{}
	copy(ip1[:], []byte("192.168.0.1"))

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

	type args struct {
		userID *ptttype.UserID_t
		passwd []byte
		ip     *ptttype.IPv4_t
	}
	tests := []struct {
		name         string
		args         args
		expectedUid  ptttype.Uid
		expectedUser *ptttype.UserecRaw
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			args:         args{userID: &userid1, passwd: []byte("123123"), ip: ip1},
			expectedUser: testSetupNewUser1,
			expectedUid:  6,
		},
		{
			args:         args{userID: &userid1, passwd: []byte("124"), ip: ip1},
			expectedUser: nil,
			wantErr:      true,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUid, gotUser, err := Login(tt.args.userID, tt.args.passwd, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUid, tt.expectedUid) {
				t.Errorf("Login() gotUid = %v, want %v", gotUid, tt.expectedUid)
			}
			if !reflect.DeepEqual(gotUser, tt.expectedUser) {
				t.Errorf("Login() gotUser = %v, want %v", gotUser, tt.expectedUser)
			}
		})
		wg.Wait()
	}
}
