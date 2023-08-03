package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestChangeEmail(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

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

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := ChangeEmail(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("ChangeEmail() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
