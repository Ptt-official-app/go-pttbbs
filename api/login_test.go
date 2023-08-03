package api

import (
	"reflect"
	"sync"
	"testing"
)

func TestLogin(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		params interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected *JwtClaim
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args: args{params: &LoginParams{
				Username: "SYSOP",
				Passwd:   "123123",
			}},
			expected: &JwtClaim{UUserID: "SYSOP"},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := Login(testIP, tt.args.params, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotResult, ok := got.(*LoginResult)
			if !ok {
				t.Errorf("Login() = %v", got)
				return
			}

			claims, _ := parseJwtClaim(gotResult.Jwt)
			if !reflect.DeepEqual(claims.UUserID, tt.expected.UUserID) {
				t.Errorf("Login() = %v claims: %v expected: %v", got, claims, tt.expected)
				return
			}
		})
		wg.Wait()
	}
}
