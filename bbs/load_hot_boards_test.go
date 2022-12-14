package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadHotBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	hbcache := []ptttype.BidInStore{9, 0, 7}
	copy(cache.Shm.Shm.HBcache[:], hbcache)
	cache.Shm.Shm.NHOTs = 3

	type args struct {
		uuserID UUserID
	}
	tests := []struct {
		name            string
		args            args
		expectedSummary []*BoardSummary
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{uuserID: "SYSOP"},
			expectedSummary: []*BoardSummary{testBoardSummary10, testBoardSummary1, testBoardSummary8},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, err := LoadHotBoards(tt.args.uuserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadHotBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for idx, each := range gotSummary {
				if idx >= len(tt.expectedSummary) {
					break
				}
				each.StatAttr = tt.expectedSummary[idx].StatAttr
				each.Total = tt.expectedSummary[idx].Total
			}

			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)
		})
	}
	wg.Wait()
}
