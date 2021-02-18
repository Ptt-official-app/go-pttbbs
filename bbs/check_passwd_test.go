package bbs

import (
	"sync"
	"testing"
)

func TestCheckPasswd(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID UUserID
		passwd  string
		ip      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{uuserID: "SYSOP", passwd: "123123", ip: "localhost"},
		},
		{
			args:    args{uuserID: "SYSOP", passwd: "123124", ip: "localhost"},
			wantErr: true,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := CheckPasswd(tt.args.uuserID, tt.args.passwd, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	wg.Wait()
}
