package ptt

import (
	"reflect"
	"sync"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/sirupsen/logrus"
)

func TestLoadGeneralBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	bsorted := [12]ptttype.BidInStore{}
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BSorted),
		unsafe.Sizeof(bsorted),
		unsafe.Pointer(&bsorted),
	)

	logrus.Infof("bsorted (by-name): %v", bsorted)
	const bsort0sz = unsafe.Sizeof(cache.Shm.Raw.BSorted[0])

	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.BSorted)+bsort0sz*uintptr(ptttype.BSORT_BY_CLASS),
		unsafe.Sizeof(bsorted),
		unsafe.Pointer(&bsorted),
	)

	logrus.Infof("bsorted (by-class): %v", bsorted)

	//move setupTest in for-loop
	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.Uid
		startIdx ptttype.SortIdx
		nBoards  int
		keyword  []byte
		bsortBy  ptttype.BSortBy
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
				bsortBy:  ptttype.BSORT_BY_NAME,
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
				bsortBy:  ptttype.BSORT_BY_NAME,
			},
			expectedSummary: []*ptttype.BoardSummaryRaw{testBoardSummary1, testBoardSummary10},
			expectedNextIdx: -1,
		},
		{
			name: "sort-by-class",
			args: args{
				user:     testUserecRaw1,
				uid:      1,
				startIdx: 0,
				nBoards:  4,
				bsortBy:  ptttype.BSORT_BY_CLASS,
			},
			expectedSummary: []*ptttype.BoardSummaryRaw{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
			expectedNextIdx: 9,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, gotNextIdx, err := LoadGeneralBoards(tt.args.user, tt.args.uid, tt.args.startIdx, tt.args.nBoards, tt.args.keyword, tt.args.bsortBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)

			if gotNextIdx != tt.expectedNextIdx {
				t.Errorf("LoadGeneralBoards() gotNextIdx = %v, want %v", gotNextIdx, tt.expectedNextIdx)
			}
		})
		wg.Wait()
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

func TestLoadHotBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	hbcache := []ptttype.BidInStore{9, 0, 7}
	cache.Shm.WriteAt(
		unsafe.Offsetof(cache.Shm.Raw.HBcache),
		unsafe.Sizeof(hbcache),
		unsafe.Pointer(&hbcache[0]),
	)
	nhots := uint8(3)
	cache.Shm.WriteAt(
		unsafe.Offsetof(cache.Shm.Raw.NHOTs),
		unsafe.Sizeof(uint8(0)),
		unsafe.Pointer(&nhots),
	)
	type args struct {
		user *ptttype.UserecRaw
		uid  ptttype.Uid
	}
	tests := []struct {
		name            string
		args            args
		expectedSummary []*ptttype.BoardSummaryRaw
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{user: testUserecRaw1, uid: 1},
			expectedSummary: []*ptttype.BoardSummaryRaw{testBoardSummary10, testBoardSummary1, testBoardSummary8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSummary, err := LoadHotBoards(tt.args.user, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadHotBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSummary, tt.expectedSummary) {
				t.Errorf("LoadHotBoards() = %v, want %v", gotSummary, tt.expectedSummary)
			}
		})
	}
}

func TestLoadBoardsByBids(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	bids := []ptttype.Bid{10, 1, 8}

	type args struct {
		user *ptttype.UserecRaw
		uid  ptttype.Uid
		bids []ptttype.Bid
	}
	tests := []struct {
		name              string
		args              args
		expectedSummaries []*ptttype.BoardSummaryRaw
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			args:              args{user: testUserecRaw1, uid: 1, bids: bids},
			expectedSummaries: []*ptttype.BoardSummaryRaw{testBoardSummary10, testBoardSummary1, testBoardSummary8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSummaries, err := LoadBoardsByBids(tt.args.user, tt.args.uid, tt.args.bids)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardsByBids() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSummaries, tt.expectedSummaries) {
				t.Errorf("LoadBoardsByBids() = %v, want %v", gotSummaries, tt.expectedSummaries)
			}
		})
	}
}
