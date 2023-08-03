package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/sirupsen/logrus"
)

func TestGetTokenInfo(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	jwt, err := CreateToken("SYSOP", "")
	if err != nil {
		logrus.Errorf("GetTokenInfo: unable to create token: jwt: %v e: %v", jwt, err)
	}
	logrus.Infof("TestGetTokenInfo: after CreateToken: jwt: %v e: %v", jwt, err)
	params0 := &GetTokenInfoParams{Jwt: jwt}
	result0 := &GetTokenInfoResult{UserID: "SYSOP"}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult interface{}
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", params: params0},
			expectedResult: result0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := GetTokenInfo(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokenInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result, _ := gotResult.(*GetTokenInfoResult)
			result.Expire = 0
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("GetTokenInfo() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
