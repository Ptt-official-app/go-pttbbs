package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestRegister(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		username string
		passwd   string
		ip       string
		email    string
		nickname []byte
		realname []byte
		career   []byte
		address  []byte
		over18   bool
	}
	tests := []struct {
		name           string
		args           args
		expectedUserID UUserID
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{username: "B1", passwd: "adsf", ip: "127.0.0.1", over18: true},
			expectedUserID: testUserec6.UUserID,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			gotUserID, err := Register(tt.args.username, tt.args.passwd, tt.args.ip, tt.args.email, tt.args.nickname, tt.args.realname, tt.args.career, tt.args.address, tt.args.over18)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "register", gotUserID, tt.expectedUserID)
		})
		wg.Wait()
	}
}
