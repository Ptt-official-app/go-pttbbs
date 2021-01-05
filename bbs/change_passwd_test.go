package bbs

import (
	"sync"
	"testing"
)

func TestChangePasswd(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		userID     UUserID
		origPasswd string
		passwd     string
		ip         string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected string
	}{
		// TODO: Add test cases.
		{
			args:     args{userID: "SYSOP", origPasswd: "non-exists", passwd: "123124", ip: "127.0.0.1"},
			wantErr:  true,
			expected: "123123",
		},
		{
			args:     args{userID: "SYSOP", origPasswd: "123123", passwd: "123124", ip: "127.0.0.1"},
			expected: "123124",
		},
		{
			args:     args{userID: "SYSOP", origPasswd: "123123", passwd: "123124", ip: "127.0.0.1"},
			wantErr:  true,
			expected: "123124",
		},
		{
			args:     args{userID: "SYSOP", origPasswd: "123124", passwd: "123123", ip: "127.0.0.1"},
			expected: "123123",
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := ChangePasswd(tt.args.userID, tt.args.origPasswd, tt.args.passwd, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("ChangePasswd() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, err := Login(string(tt.args.userID), tt.expected, tt.args.ip)
			if err != nil {
				t.Errorf("ChangePasswd() unable to login: userID: %v expected: %v ip: %v", tt.args.userID, tt.expected, tt.args.ip)
			}
		})
		wg.Wait()
	}
}
