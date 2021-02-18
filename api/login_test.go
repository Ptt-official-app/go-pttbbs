package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"gopkg.in/square/go-jose.v2/jwt"
)

func TestLogin(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		params interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args: args{params: &LoginParams{
				Username: "SYSOP",
				Passwd:   "123123",
			}},
			expected: &JwtClaim{UUserID: bbs.UUserID("SYSOP")},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := Login(testIP, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotResult, ok := got.(*LoginResult)
			if !ok {
				t.Errorf("Login() = %v", got)
				return
			}

			claims := &JwtClaim{}
			token, _ := jwt.ParseSigned(gotResult.Jwt)
			_ = token.UnsafeClaimsWithoutVerification(claims)
			wantJwt, _ := tt.expected.(*JwtClaim)
			if !reflect.DeepEqual(claims.UUserID, wantJwt.UUserID) {
				t.Errorf("Login() = %v claims: %v expected: %v", got, claims, tt.expected)
				return
			}
		})
		wg.Wait()
	}
}
