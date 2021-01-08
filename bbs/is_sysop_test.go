package bbs

import (
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsValid := IsSysop(tt.args.uuserID, tt.args.perm); gotIsValid != tt.expectedIsValid {
				t.Errorf("IsSysop() = %v, want %v", gotIsValid, tt.expectedIsValid)
			}
		})
	}
}
