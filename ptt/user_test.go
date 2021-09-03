package ptt

import (
	"bytes"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/sirupsen/logrus"
)

func Test_killUser(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("CodingMan"))

	type args struct {
		uid    ptttype.UID
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{uid: 1, userID: userID1},
		},
		{
			args: args{uid: 1, userID: userID1},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			if err := killUser(tt.args.uid, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("killUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			user, err := passwdSyncQuery(tt.args.uid)
			if err != nil {
				t.Errorf("killUser: unable to query: e: %v", err)
			}

			if !bytes.Equal(user.UserID[:], ptttype.EMPTY_USER_ID[:]) {
				t.Errorf("killUser: unable to kill: userID: %v", string(user.UserID[:]))
			}
		})
		wg.Wait()
	}
}

func Test_tryDeleteHomePath(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	defer os.RemoveAll("./testcase/tmp")

	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("CodingMan"))

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{userID1},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			homepath := path.SetHomePath(tt.args.userID)
			_, err := os.Stat(homepath)
			if err != nil {
				t.Errorf("tryDeleteHomePath: home-path not exists: homepath: %v", homepath)
			}

			if err := tryDeleteHomePath(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("tryDeleteHomePath() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, err = os.Stat(homepath)
			if err == nil {
				t.Errorf("tryDeleteHomePath: still with hoem-path: homepath: %v", homepath)
			}
		})
		wg.Wait()
	}
}

func TestChangePasswd(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("SYSOP"))
	origPasswd0 := []byte("123123")
	passwd0 := []byte("123124")
	badPasswd0 := []byte("non-exists")

	expected0 := origPasswd0
	expected1 := passwd0
	expected2 := origPasswd0

	type args struct {
		userID     *ptttype.UserID_t
		origPasswd []byte
		passwd     []byte
	}
	tests := []struct {
		name     string
		args     args
		expected []byte
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{userID: userID0, origPasswd: badPasswd0, passwd: passwd0},
			wantErr:  true,
			expected: expected0,
		},
		{
			args:     args{userID: userID0, origPasswd: origPasswd0, passwd: passwd0},
			expected: expected1,
		},
		{
			args:     args{userID: userID0, origPasswd: passwd0, passwd: origPasswd0},
			expected: expected2,
		},
	}

	ip := &ptttype.IPv4_t{}
	copy(ip[:], []byte("127.0.0.1"))
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			err := ChangePasswd(tt.args.userID, tt.args.origPasswd, tt.args.passwd, ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePasswd() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, _, err = LoginQuery(tt.args.userID, tt.expected, ip)
			if err != nil {
				t.Errorf("ChangePasswd: unable to login: e: %v", err)
			}
		})
		wg.Wait()
	}
}

func TestChangeEmail(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("SYSOP"))

	email0 := &ptttype.Email_t{}
	copy(email0[:], "test@ptt.test")

	type args struct {
		userID *ptttype.UserID_t
		email  *ptttype.Email_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				userID: userID0,
				email:  email0,
			},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := ChangeEmail(tt.args.userID, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("ChangeEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			user, _ := GetUser(tt.args.userID)
			testutil.TDeepEqual(t, "email", &user.Email, tt.args.email)
		})
		wg.Wait()
	}
}

func TestCheckPasswd(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("SYSOP"))

	ip := &ptttype.IPv4_t{}
	copy(ip[:], []byte("localhost"))

	passwd0 := []byte("123123")

	passwd1 := []byte("123124")

	type args struct {
		userID *ptttype.UserID_t
		passwd []byte
		ip     *ptttype.IPv4_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{userID: userID0, passwd: passwd0, ip: ip},
		},
		{
			args:    args{userID: userID0, passwd: passwd1, ip: ip},
			wantErr: true,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := CheckPasswd(tt.args.userID, tt.args.passwd, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}

func TestChangeUserLevel2(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		userID *ptttype.UserID_t
		perm   ptttype.PERM2
		isSet  bool
	}
	tests := []struct {
		name               string
		args               args
		expectedUserLevel2 ptttype.PERM2
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				userID: &testUserecRaw1.UserID,
				perm:   ptttype.PERM2_ID_EMAIL,
				isSet:  true,
			},
			expectedUserLevel2: ptttype.PERM2_ID_EMAIL,
		},
		{
			args: args{
				userID: &testUserecRaw1.UserID,
				perm:   ptttype.PERM2_ID_EMAIL,
				isSet:  false,
			},
			expectedUserLevel2: ptttype.PERM2_INVALID,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUserLevel2, err := ChangeUserLevel2(tt.args.userID, tt.args.perm, tt.args.isSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeUserLevel2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserLevel2, tt.expectedUserLevel2) {
				t.Errorf("ChangeUserLevel2() = %v, want %v", gotUserLevel2, tt.expectedUserLevel2)
			}
		})
		wg.Wait()
	}
}

func TestGetUid(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("SYSOP"))

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name        string
		args        args
		expectedUid ptttype.UID
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{userID: userID0},
			expectedUid: 1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUid, err := GetUID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUid, tt.expectedUid) {
				t.Errorf("GetUid() = %v, want %v", gotUid, tt.expectedUid)
			}
		})
		wg.Wait()
	}
}

func TestSetUserPerm(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	origPerm := testUserecRaw2.UserLevel
	newPerm := ptttype.PERM_DEFAULT | ptttype.PERM_ADMIN | ptttype.PERM_LOGINOK

	logrus.Infof("TestSetUserPerm: origPerm: %v newPerm: %v", origPerm, newPerm)

	type args struct {
		userec    *ptttype.UserecRaw
		setUid    ptttype.UID
		setUserec *ptttype.UserecRaw
		perm      ptttype.PERM
	}
	tests := []struct {
		name             string
		args             args
		expectedNewPerm  ptttype.PERM
		expectedNewPerm2 ptttype.PERM
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			args:             args{userec: testUserecRaw1, setUid: 2, setUserec: testUserecRaw2, perm: newPerm},
			expectedNewPerm:  newPerm,
			expectedNewPerm2: newPerm,
		},
		{
			args:             args{userec: testUserecRaw1, setUid: 2, setUserec: testUserecRaw2, perm: origPerm},
			expectedNewPerm:  origPerm,
			expectedNewPerm2: origPerm,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotNewPerm, err := SetUserPerm(tt.args.userec, tt.args.setUid, tt.args.setUserec, tt.args.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUserPerm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewPerm, tt.expectedNewPerm) {
				t.Errorf("SetUserPerm() = %v, want %v", gotNewPerm, tt.expectedNewPerm)
			}

			newUser, _ := InitCurrentUserByUID(tt.args.setUid)
			if !reflect.DeepEqual(newUser.UserLevel, tt.expectedNewPerm) {
				t.Errorf("SetUserPerm() newUser: %v want: %v", newUser.UserLevel, tt.expectedNewPerm)
			}
		})
		wg.Wait()
	}
}
