package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadBoardSummary(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID  UUserID
		bboardID BBoardID
	}
	tests := []struct {
		name            string
		args            args
		expectedSummary *BoardSummary
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{uuserID: "SYSOP", bboardID: "6_ALLPOST"},
			expectedSummary: testBoardSummary6,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, err := LoadBoardSummary(tt.args.uuserID, tt.args.bboardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotSummary, tt.expectedSummary)
		})
		wg.Wait()
	}
}
