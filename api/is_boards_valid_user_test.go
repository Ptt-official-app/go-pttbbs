package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func TestIsBoardsValidUser(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	_ = ptt.SetupNewUser(testNewPostUserRaw1)

	params0 := &IsBoardsValidUserParams{
		BoardIDs: []bbs.BBoardID{"1_SYSOP", "7_deleted", "2_1...........", "5_2..........."},
	}
	expected0 := &IsBoardsValidUserResult{
		IsValid: map[bbs.BBoardID]bool{"1_SYSOP": true, "7_deleted": true, "2_1...........": true, "5_2...........": true},
	}
	expected1 := &IsBoardsValidUserResult{
		IsValid: map[bbs.BBoardID]bool{"1_SYSOP": true, "7_deleted": false, "2_1...........": false, "5_2...........": true},
	}
	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult *IsBoardsValidUserResult
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", params: params0},
			expectedResult: expected0,
		},
		{
			args:           args{remoteAddr: testIP, uuserID: "A1", params: params0},
			expectedResult: expected1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := IsBoardsValidUser(tt.args.remoteAddr, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsBoardsValidUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("IsBoardsValidUser() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
