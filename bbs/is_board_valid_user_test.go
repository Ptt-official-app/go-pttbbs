package bbs

import (
	"sync"
	"testing"
)

func TestIsBoardValidUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID UUserID
		boardID BBoardID
	}
	tests := []struct {
		name            string
		args            args
		expectedIsValid bool
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{uuserID: "SYSOP", boardID: "10_WhoAmI"},
			expectedIsValid: true,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotIsValid, err := IsBoardValidUser(tt.args.uuserID, tt.args.boardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsBoardValidUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsValid != tt.expectedIsValid {
				t.Errorf("IsBoardValidUser() = %v, want %v", gotIsValid, tt.expectedIsValid)
			}
		})
	}
	wg.Wait()
}
