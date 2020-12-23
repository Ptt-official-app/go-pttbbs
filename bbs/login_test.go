package bbs

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestLogin(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		username string
		passwd   string
		ip       string
	}
	tests := []struct {
		name     string
		args     args
		expected *Userec
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{"SYSOP", "123123", "127.0.0.1"},
			expected: testUserec1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Login(tt.args.username, tt.args.passwd, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			types.TDeepEqual(t, "login", got, tt.expected)
		})
	}
}
