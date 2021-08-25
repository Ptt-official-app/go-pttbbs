package bbs

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadFullClassBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID  UUserID
		startBid ptttype.Bid
		nBoards  int
	}
	tests := []struct {
		name              string
		args              args
		expectedSummaries []*BoardSummary
		expectedNextBid   ptttype.Bid
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			args:              args{uuserID: "SYSOP", startBid: 1, nBoards: 100},
			expectedSummaries: []*BoardSummary{testClassSummary2, testClassSummary5},
		},
		{
			args:              args{uuserID: "SYSOP", startBid: 1, nBoards: 1},
			expectedSummaries: []*BoardSummary{testClassSummary2},
			expectedNextBid:   5,
		},
		{
			args:              args{uuserID: "SYSOP", startBid: 3, nBoards: 1},
			expectedSummaries: []*BoardSummary{testClassSummary5},
		},
		{
			args:              args{uuserID: "SYSOP", startBid: 2, nBoards: 1},
			expectedSummaries: []*BoardSummary{testClassSummary2},
			expectedNextBid:   5,
		},
		{
			args:              args{uuserID: "SYSOP", startBid: 6, nBoards: 1},
			expectedSummaries: []*BoardSummary{},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, gotNextBid, err := LoadFullClassBoards(tt.args.uuserID, tt.args.startBid, tt.args.nBoards)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFullClassBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotSummaries, tt.expectedSummaries)
			if !reflect.DeepEqual(gotNextBid, tt.expectedNextBid) {
				t.Errorf("LoadFullClassBoards() gotNextBid = %v, want %v", gotNextBid, tt.expectedNextBid)
			}
		})
		wg.Wait()
	}
}
