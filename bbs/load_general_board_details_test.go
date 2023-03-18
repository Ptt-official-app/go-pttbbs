package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadGeneralBoardDetails(t *testing.T) {
	setupTest()
	defer teardownTest()
	type args struct {
		uuserID     UUserID
		startIdxStr string
		nBoards     int
		isAsc       bool
		bsortBy     ptttype.BSortBy
	}
	tests := []struct {
		name           string
		args           args
		wantDetails    []*BoardDetail
		wantNextIdxStr string
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{uuserID: "SYSOP3", startIdxStr: "", nBoards: 4, isAsc: true, bsortBy: ptttype.BSORT_BY_NAME},
			wantDetails:    []*BoardDetail{testClassDetail2, testClassDetail5, testBoardDetail12, testBoardDetail6},
			wantNextIdxStr: "deleted",
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotDetails, gotNextIdxStr, err := LoadGeneralBoardDetails(tt.args.uuserID, tt.args.startIdxStr, tt.args.nBoards, tt.args.isAsc, tt.args.bsortBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoardDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, each := range gotDetails {
				each.Total = 0
				each.LastPostTime = 0
			}
			testutil.TDeepEqual(t, "details", gotDetails, tt.wantDetails)
			if gotNextIdxStr != tt.wantNextIdxStr {
				t.Errorf("LoadGeneralBoardDetails() gotNextIdxStr = %v, want %v", gotNextIdxStr, tt.wantNextIdxStr)
			}
		})
	}
	wg.Wait()
}
