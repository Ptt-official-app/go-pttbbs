package api

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadBoardsByBids(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	params := &LoadBoardsByBidsParams{
		Bids: []ptttype.Bid{6, 7, 11, 8},
	}

	expected := &LoadBoardsByBidsResult{
		Boards: []*bbs.BoardSummary{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult interface{}
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", params: params},
			expectedResult: expected,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := LoadBoardsByBids(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardsByBids() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "result", gotResult, tt.expectedResult)
		})
		wg.Wait()
	}
}
