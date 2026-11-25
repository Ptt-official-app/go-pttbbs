package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestNewBoardDetailFromRaw(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		boardDetailRaw *ptttype.BoardDetailRaw
		bid            ptttype.Bid
	}
	tests := []struct {
		name     string
		args     args
		expected *BoardDetail
	}{
		// TODO: Add test cases.
		{
			args:     args{testBoardDetailRaw3, 1},
			expected: testBoardDetail3,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)

		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got := NewBoardDetailFromRaw(tt.args.boardDetailRaw, tt.args.bid)
			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
		wg.Wait()
	}
}

func TestNewBoardDetailFromRaw_limitLogins(t *testing.T) {
	setupTest()
	defer teardownTest()

	raw := *testBoardDetailRaw3
	raw.VoteLimitLogins = 7
	raw.PostLimitLogins = 9

	got := NewBoardDetailFromRaw(&raw, 1)
	if got.VoteLimitLogins != 7 {
		t.Errorf("VoteLimitLogins: got: %v expected: 7", got.VoteLimitLogins)
	}
	if got.PostLimitLogins != 9 {
		t.Errorf("PostLimitLogins: got: %v expected: 9", got.PostLimitLogins)
	}
}
