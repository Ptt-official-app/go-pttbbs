package ptt

import (
	"bytes"
	"os"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func Test_killUser(t *testing.T) {
	setupTest()
	defer teardownTest()

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
	setupTest()
	defer func() {
		teardownTest()
		os.RemoveAll("./testcase/tmp")
	}()

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

func Test_hasUserPerm(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		user *ptttype.UserecRaw
		perm ptttype.PERM
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{user: testNewRegister1, perm: ptttype.PERM_POST},
			expected: false,
		},
		{
			args:     args{user: testNewRegister1, perm: ptttype.PERM_BASIC},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasUserPerm(tt.args.user, tt.args.perm); got != tt.expected {
				t.Errorf("hasUserPerm() = %v, want %v", got, tt.expected)
			}
		})
	}
}
