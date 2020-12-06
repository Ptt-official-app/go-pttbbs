package api

import (
	"reflect"
	"testing"

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
				UserID: "SYSOP",
				Passwd: "123123",
			}},
			expected: &JwtClaim{UserID: "SYSOP"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Login(tt.args.params)
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
			if !reflect.DeepEqual(claims.UserID, wantJwt.UserID) {
				t.Errorf("Login() = %v claims: %v expected: %v", got, claims, tt.expected)
				return
			}
		})
	}
}
