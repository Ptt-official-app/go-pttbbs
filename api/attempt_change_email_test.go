package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/sirupsen/logrus"
)

func TestAttemptChangeEmail(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	params0 := &AttemptChangeEmailParams{
		ClientInfo: "test_clientinfo",
		Passwd:     "123123",
		Email:      "test@ptt.test",
	}

	path0 := &AttemptChangeEmailPath{
		UserID: "SYSOP",
	}

	result0 := &AttemptChangeEmailResult{
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
		expectedResult *AttemptChangeEmailResult
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", params: params0, path: path0},
			expectedResult: result0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := AttemptChangeEmail(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttemptChangeEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			result := gotResult.(*AttemptChangeEmailResult)
			logrus.Infof("AttemptChangeEmail: result: %v", result)
			tt.expectedResult.Jwt = result.Jwt
			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("AttemptChangeEmail() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
