package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadClassBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID UUserID
		clsBid  ptttype.Bid
		bsortBy ptttype.BSortBy
	}
	tests := []struct {
		name              string
		args              args
		expectedSummaries []*BoardSummary
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			args:              args{uuserID: "SYSOP", clsBid: 1, bsortBy: ptttype.BSORT_BY_NAME},
			expectedSummaries: []*BoardSummary{testClassSummary2, testClassSummary5},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, err := LoadClassBoards(tt.args.uuserID, tt.args.clsBid, tt.args.bsortBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadClassBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotSummaries, tt.expectedSummaries)
		})
		wg.Wait()
	}
}
