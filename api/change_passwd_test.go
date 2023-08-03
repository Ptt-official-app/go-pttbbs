package api

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestChangePasswd(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	params0 := &ChangePasswdParams{
		ClientInfo: "test_clientinfo",
		OrigPasswd: "123123",
		Passwd:     "123124",
	}
	path0 := &ChangePasswdPath{
		UserID: "SYSOP",
	}
	result0 := &ChangePasswdResult{
		UserID:    "SYSOP",
		TokenType: "bearer",
	}

	params1 := &ChangePasswdParams{
		ClientInfo: "test_clientinfo",
		OrigPasswd: "123123",
		Passwd:     "123124",
	}

	params2 := &ChangePasswdParams{
		ClientInfo: "test_clientinfo",
		OrigPasswd: "123124",
		Passwd:     "123123",
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
		expectedResult *ChangePasswdResult
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: "127.0.0.1", uuserID: bbs.UUserID("SYSOP"), params: params0, path: path0},
			expectedResult: result0,
		},
		{
			args:           args{remoteAddr: "127.0.0.1", uuserID: bbs.UUserID("SYSOP"), params: params1, path: path0},
			expectedResult: nil,
			wantErr:        true,
		},
		{
			args:           args{remoteAddr: "127.0.0.1", uuserID: bbs.UUserID("SYSOP"), params: params2, path: path0},
			expectedResult: result0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := ChangePasswd(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePasswd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			result := gotResult.(*ChangePasswdResult)
			tt.expectedResult.Jwt = result.Jwt
			testutil.TDeepEqual(t, "result", result, tt.expectedResult)
		})
		wg.Wait()
	}
}
