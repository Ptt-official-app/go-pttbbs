package api

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func Test_userIsValidUser(t *testing.T) {
	setupTest()
	defer teardownTest()

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

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if gotIsValid := userInfoIsValidUser(tt.args.uuserID, tt.args.queryUUserID); gotIsValid != tt.expectedIsValid {
				t.Errorf("userIsValidUser() = %v, want %v", gotIsValid, tt.expectedIsValid)
			}
		})
	}
	wg.Wait()
}
