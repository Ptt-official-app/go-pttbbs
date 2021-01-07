package api

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestChangeEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	//until Fri Nov 14 01:28:37 EST 2245
	params0 := &ChangeEmailParams{
		Jwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGkiOiJ0ZXN0X2NsaWVudGluZm8iLCJlbWwiOiJ0ZXN0QHB0dC50ZXN0IiwiZXhwIjo4NzA1NjAwOTE3LCJzdWIiOiJTWVNPUCJ9.IPa_yF9BYbnnrjwguwPnE7fxpV13bhfgex74-ok-VzE",
	}
	path0 := &ChangeEmailPath{
		UserID: "SYSOP",
	}
	result0 := &ChangeEmailResult{
		UserID: "SYSOP",
		Email:  "test@ptt.test",
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
		path       interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult interface{}
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", params: params0, path: path0},
			expectedResult: result0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := ChangeEmail(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("ChangeEmail() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
	}
}
