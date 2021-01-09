package api

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestGetTokenInfo(t *testing.T) {

	jwt, _ := CreateToken("SYSOP", "")
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := GetTokenInfo(tt.args.remoteAddr, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokenInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("GetTokenInfo() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
	}
}
