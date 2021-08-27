package bbs

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func TestIsBoardsValidUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testNewPostUserRaw1)

	type args struct {
		uuserID  UUserID
		boardIDs []BBoardID
	}
	tests := []struct {
		name            string
		args            args
		expectedIsValid map[BBoardID]bool
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{uuserID: "SYSOP", boardIDs: []BBoardID{"1_SYSOP", "10_WhoAmI", "7_deleted", "2_1..........."}},
			expectedIsValid: map[BBoardID]bool{"1_SYSOP": true, "10_WhoAmI": true, "7_deleted": true, "2_1...........": true},
		},
		{
			args:            args{uuserID: "A1", boardIDs: []BBoardID{"1_SYSOP", "10_WhoAmI", "7_deleted", "2_1...........", "5_2..........."}},
			expectedIsValid: map[BBoardID]bool{"1_SYSOP": true, "10_WhoAmI": true, "7_deleted": false, "2_1...........": false, "5_2...........": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsValid, err := IsBoardsValidUser(tt.args.uuserID, tt.args.boardIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsBoardsValidUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIsValid, tt.expectedIsValid) {
				t.Errorf("IsBoardsValidUser() = %v, want %v", gotIsValid, tt.expectedIsValid)
			}
		})
	}
}
