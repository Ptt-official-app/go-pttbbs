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
)

func Test_killUser(t *testing.T) {
	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("CodingMan"))

	type args struct {
		uid    ptttype.Uid
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

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
	}
}

func Test_tryDeleteHomePath(t *testing.T) {

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer func() {
				teardownTest()
				os.RemoveAll("./testcase/tmp")
			}()

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
	}
}

func TestChangePasswd(t *testing.T) {
	setupTest()
	defer teardownTest()

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
	setupTest()
	defer teardownTest()

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeEmail(tt.args.userID, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("ChangeEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			user, _ := GetUser(tt.args.userID)
			testutil.TDeepEqual(t, "email", &user.Email, tt.args.email)
		})
	}
}

func TestCheckPasswd(t *testing.T) {
	setupTest()
	defer teardownTest()

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPasswd(tt.args.userID, tt.args.passwd, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangeUserLevel2(t *testing.T) {
	setupTest()
	defer teardownTest()

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
