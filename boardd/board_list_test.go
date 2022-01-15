package boardd

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadHotBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		user *ptttype.UserecRaw
		uid  ptttype.UID
	}
	tests := []struct {
		name            string
		args            args
		expectedSummary []*ptttype.BoardSummaryRaw
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{},
			expectedSummary: []*ptttype.BoardSummaryRaw{testBoardSummary10, testBoardSummary1, testBoardSummary8},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, err := LoadHotBoards(tt.args.user, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadHotBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotSummary, tt.expectedSummary)
		})
		wg.Wait()
	}
}
