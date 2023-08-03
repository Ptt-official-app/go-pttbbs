package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

func TestGetRefreshTokenInfo(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	refreshJwt, _ := CreateRefreshToken("SYSOP", "")

	params0 := &GetRefreshTokenInfoParams{
		Jwt: refreshJwt,
	}
	expected0 := &GetRefreshTokenInfoResult{
		UserID: "SYSOP",
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
		c          *gin.Context
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{uuserID: "SYSOP", params: params0},
			expected: expected0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := GetRefreshTokenInfo(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRefreshTokenInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			result, _ := gotResult.(*GetRefreshTokenInfoResult)
			result.Expire = 0
			if !reflect.DeepEqual(gotResult, tt.expected) {
				t.Errorf("GetRefreshTokenInfo() = %v, want %v", gotResult, tt.expected)
			}
		})
		wg.Wait()
	}
}
