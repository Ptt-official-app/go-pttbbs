package ptt

import (
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	log "github.com/sirupsen/logrus"
)

func TestLoadGeneralBoards(t *testing.T) {
	//move setupTest in for-loop
	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.Uid
		startIdx ptttype.SortIdx
		nBoards  int
		keyword  []byte
	}
	tests := []struct {
		name            string
		args            args
		expectedSummary []*ptttype.BoardSummaryRaw
		expectedNextIdx ptttype.SortIdx
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:     testUserecRaw1,
				uid:      1,
				startIdx: 0,
				nBoards:  4,
			},
			expectedSummary: []*ptttype.BoardSummaryRaw{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
			expectedNextIdx: 8,
		},
		{
			args: args{
				user:     testUserecRaw1,
				uid:      1,
				startIdx: 10,
				nBoards:  4,
			},
			expectedSummary: []*ptttype.BoardSummaryRaw{testBoardSummary1, testBoardSummary10},
			expectedNextIdx: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			cache.ReloadBCache()

			bsorted := [12]ptttype.BidInStore{}
			cache.Shm.ReadAt(
				unsafe.Offsetof(cache.Shm.Raw.BSorted),
				unsafe.Sizeof(bsorted),
				unsafe.Pointer(&bsorted),
			)

			log.Infof("bsorted: %v", bsorted)

			gotSummary, gotNextIdx, err := LoadGeneralBoards(tt.args.user, tt.args.uid, tt.args.startIdx, tt.args.nBoards, tt.args.keyword)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)

			if gotNextIdx != tt.expectedNextIdx {
				t.Errorf("LoadGeneralBoards() gotNextIdx = %v, want %v", gotNextIdx, tt.expectedNextIdx)
			}
		})
	}
}

func TestLoadBoardSummary(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	bsorted := [12]ptttype.BidInStore{}
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BSorted),
		unsafe.Sizeof(bsorted),
		unsafe.Pointer(&bsorted),
	)

	type args struct {
		user *ptttype.UserecRaw
		uid  ptttype.Uid
		bid  ptttype.Bid
	}
	tests := []struct {
		name            string
		args            args
		expectedSummary *ptttype.BoardSummaryRaw
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{user: testUserecRaw1, uid: 1, bid: 10},
			expectedSummary: testBoardSummary10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSummary, err := LoadBoardSummary(tt.args.user, tt.args.uid, tt.args.bid)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)
		})
	}
}
