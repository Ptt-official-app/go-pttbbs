package api

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func Test_userIsValidUser(t *testing.T) {
	type args struct {
		uuserID      bbs.UUserID
		queryUUserID bbs.UUserID
	}
	tests := []struct {
		name            string
		args            args
		expectedIsValid bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsValid := userInfoIsValidUser(tt.args.uuserID, tt.args.queryUUserID); gotIsValid != tt.expectedIsValid {
				t.Errorf("userIsValidUser() = %v, want %v", gotIsValid, tt.expectedIsValid)
			}
		})
	}
}
