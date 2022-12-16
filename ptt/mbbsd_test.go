package ptt

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func TestLoginQuery(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

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
		expectedUid ptttype.UID
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
			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
		wg.Wait()
	}
}

func Test_newUserInfoRaw(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		uid  ptttype.UID
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
	setupTest(t.Name())
	defer teardownTest(t.Name())

	SetupNewUser(testSetupNewUser1)

	userid1 := ptttype.UserID_t{}
	copy(userid1[:], []byte("A0"))

	ip1 := &ptttype.IPv4_t{}
	copy(ip1[:], []byte("192.168.0.1"))

	cache.Shm.Shm.CurrSorted = 0

	cache.Shm.Shm.UTMPNumber = 5

	sorted := []ptttype.UtmpID{2, 1, 0, 4, 3}
	copy(cache.Shm.Shm.Sorted[0][ptttype.SORT_BY_PID][:], sorted)

	uinfo := []ptttype.UserInfoRaw{testUserInfo1, testUserInfo2, testUserInfo3, testUserInfo4, testUserInfo5, testUserInfo6}
	copy(cache.Shm.Shm.UInfo[:], uinfo)

	// setup login data
	testSetupNewUser1.LastHost = *ip1
	testSetupNewUser1.NumLoginDays = testSetupNewUser1.NumLoginDays + 1
	testSetupNewUser1.LastLogin = types.NowTS()
	testSetupNewUser1.LastSeen = types.NowTS()

	type args struct {
		userID *ptttype.UserID_t
		passwd []byte
		ip     *ptttype.IPv4_t
	}
	tests := []struct {
		name         string
		args         args
		expectedUid  ptttype.UID
		expectedUser *ptttype.UserecRaw
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			args:         args{userID: &userid1, passwd: []byte("123123"), ip: ip1},
			expectedUser: testSetupNewUser1,
			expectedUid:  41,
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
			testutil.TDeepEqual(t, "user", gotUser, tt.expectedUser)
		})
		wg.Wait()
	}
}
