package bbs

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestRegister(t *testing.T) {
	type args struct {
		username string
		passwd   string
		ip       string
		email    string
		nickname string
		realname string
		career   string
		address  string
		over18   bool
	}
	tests := []struct {
		name         string
		args         args
		expectedUser *Userec
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			args:         args{username: "B1", passwd: "adsf", ip: "127.0.0.1", over18: true},
			expectedUser: testUserec6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			gotUser, err := Register(tt.args.username, tt.args.passwd, tt.args.ip, tt.args.email, tt.args.nickname, tt.args.realname, tt.args.career, tt.args.address, tt.args.over18)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotUser.Firstlogin = 0
			gotUser.Lastlogin = 0

			types.TDeepEqual(t, gotUser, tt.expectedUser)
		})
	}
}
