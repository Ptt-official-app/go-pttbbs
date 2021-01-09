package api

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestChangeEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	jwt, _ := CreateEmailToken("SYSOP", "", "test@ptt.test", CONTEXT_CHANGE_EMAIL)

	params0 := &ChangeEmailParams{
		Jwt: jwt,
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
