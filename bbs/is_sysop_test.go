package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestIsSysop(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID UUserID
		perm    ptttype.PERM
	}
	tests := []struct {
		name            string
		args            args
		expectedIsValid bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				uuserID: "SYSOP",
				perm:    ptttype.PERM_SYSSUBOP | ptttype.PERM_BM,
			},
			expectedIsValid: true,
		},
		{
			args: args{
				uuserID: "SYSOP",
				perm:    ptttype.PERM_ACCOUNTS,
			},
			expectedIsValid: false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if gotIsValid := IsSysop(tt.args.uuserID, tt.args.perm); gotIsValid != tt.expectedIsValid {
				t.Errorf("IsSysop() = %v, want %v", gotIsValid, tt.expectedIsValid)
			}
		})
	}
	wg.Wait()
}
