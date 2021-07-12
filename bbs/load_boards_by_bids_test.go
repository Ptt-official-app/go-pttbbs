package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadBoardsByBids(t *testing.T) {
	setupTest()
	defer teardownTest()

	bids := []ptttype.Bid{6, 7, 11, 8}

	type args struct {
		uuserID UUserID
		bids    []ptttype.Bid
	}
	tests := []struct {
		name              string
		args              args
		expectedSummaries []*BoardSummary
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			args:              args{uuserID: "SYSOP", bids: bids},
			expectedSummaries: []*BoardSummary{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, err := LoadBoardsByBids(tt.args.uuserID, tt.args.bids)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardsByBids() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "summaries", gotSummaries, tt.expectedSummaries)
		})
		wg.Wait()
	}
}
