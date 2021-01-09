package api

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestAttemptSetIDEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	params0 := &AttemptSetIDEmailParams{
		ClientInfo: "test_clientinfo",
		Passwd:     "123123",
		Email:      "test@ptt.test",
	}

	path0 := &AttemptSetIDEmailPath{
		UserID: "SYSOP",
	}

	result0 := &AttemptSetIDEmailResult{
		UserID: "SYSOP",
		Jwt:    "",
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
		expectedResult *AttemptSetIDEmailResult
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
			gotResult, err := AttemptSetIDEmail(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttemptSetIDEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result := gotResult.(*AttemptSetIDEmailResult)
			tt.expectedResult.Jwt = result.Jwt
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("AttemptSetIDEmail() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
	}
}
