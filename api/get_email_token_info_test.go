package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestGetEmailTokenInfo(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	jwt, _ := CreateEmailToken("SYSOP", "", "test@ptt.test", CONTEXT_CHANGE_EMAIL)
	params0 := &GetEmailTokenInfoParams{Jwt: jwt, Context: CONTEXT_CHANGE_EMAIL}
	result0 := &GetEmailTokenInfoResult{UserID: "SYSOP", Email: "test@ptt.test"}

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
			gotResult, err := GetEmailTokenInfo(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEmailTokenInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			result, _ := gotResult.(*GetEmailTokenInfoResult)
			result.Expire = 0
			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("GetEmailTokenInfo() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
