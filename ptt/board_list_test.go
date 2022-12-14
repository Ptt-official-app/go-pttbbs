package ptt

import (
	"reflect"
	"sync"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/sirupsen/logrus"
)

func TestLoadGeneralBoards(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	bsorted := &cache.Shm.Shm.BSorted[ptttype.BSORT_BY_NAME]
	logrus.Infof("bsorted (by-name): %v", bsorted)

	const bsort0sz = unsafe.Sizeof(cache.Shm.Raw.BSorted[0])

	bsorted = &cache.Shm.Shm.BSorted[ptttype.BSORT_BY_CLASS]
	logrus.Infof("bsorted (by-class): %v", bsorted)

	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.UID
		startIdx ptttype.SortIdx
		nBoards  int
		title    []byte
		keyword  []byte
		isAsc    bool
		bsortBy  ptttype.BSortBy
	}
	tests := []struct {
		name                string
		args                args
		expectedSummaries   []*ptttype.BoardSummaryRaw
		expectedNextSummary *ptttype.BoardSummaryRaw
		wantErr             bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:     testUserecRaw4,
				uid:      11,
				startIdx: 1,
				nBoards:  4,
				bsortBy:  ptttype.BSORT_BY_NAME,
				isAsc:    true,
			},
			expectedSummaries:   []*ptttype.BoardSummaryRaw{testBoardSummary6, testBoardSummary11, testBoardSummary8, testBoardSummary9},
			expectedNextSummary: testBoardSummary1,
		},
		{
			args: args{
				user:     testUserecRaw4,
				uid:      11,
				startIdx: 5,
				nBoards:  4,
				bsortBy:  ptttype.BSORT_BY_NAME,
				isAsc:    false,
			},
			expectedSummaries:   []*ptttype.BoardSummaryRaw{testBoardSummary6},
			expectedNextSummary: nil,
		},
		{
			args: args{
				user:     testUserecRaw4,
				uid:      11,
				startIdx: 10,
				nBoards:  4,
				bsortBy:  ptttype.BSORT_BY_NAME,
				isAsc:    true,
			},
			expectedSummaries:   []*ptttype.BoardSummaryRaw{testBoardSummary1, testBoardSummary10},
			expectedNextSummary: nil,
		},
		{
			args: args{
				user:     testUserecRaw4,
				uid:      11,
				startIdx: 1,
				nBoards:  4,
				bsortBy:  ptttype.BSORT_BY_NAME,
				isAsc:    true,
				title:    []byte{'o'},
			},
			expectedSummaries:   []*ptttype.BoardSummaryRaw{testBoardSummary6, testBoardSummary8, testBoardSummary10},
			expectedNextSummary: nil,
		},
		{
			args: args{
				user:     testUserecRaw4,
				uid:      11,
				startIdx: 1,
				nBoards:  4,
				bsortBy:  ptttype.BSORT_BY_NAME,
				isAsc:    true,
				keyword:  []byte{'o'},
			},
			expectedSummaries:   []*ptttype.BoardSummaryRaw{testBoardSummary6, testBoardSummary8, testBoardSummary9, testBoardSummary1},
			expectedNextSummary: testBoardSummary10,
		},
		{
			name: "sort-by-class",
			args: args{
				user:     testUserecRaw4,
				uid:      11,
				startIdx: 1,
				nBoards:  4,
				bsortBy:  ptttype.BSORT_BY_CLASS,
				isAsc:    true,
			},
			expectedSummaries:   []*ptttype.BoardSummaryRaw{testBoardSummary6, testBoardSummary11, testBoardSummary8, testBoardSummary9},
			expectedNextSummary: testBoardSummary1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, gotNextSummary, err := LoadGeneralBoards(tt.args.user, tt.args.uid, tt.args.startIdx, tt.args.nBoards, tt.args.title, tt.args.keyword, tt.args.isAsc, tt.args.bsortBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, each := range gotSummaries {
				each.Total = 0
				each.LastPostTime = 0
			}
			testutil.TDeepEqual(t, "summaries", gotSummaries, tt.expectedSummaries)

			if gotNextSummary != nil {
				gotNextSummary.Total = 0
				gotNextSummary.LastPostTime = 0
			}
			testutil.TDeepEqual(t, "nextSummary", gotNextSummary, tt.expectedNextSummary)
		})
		wg.Wait()
	}
}

func TestLoadBoardSummary(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	type args struct {
		user *ptttype.UserecRaw
		uid  ptttype.UID
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

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, err := LoadBoardSummary(tt.args.user, tt.args.uid, tt.args.bid)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotSummary.Total = 0
			gotSummary.LastPostTime = 0
			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)
		})
		wg.Wait()
	}
}

func TestLoadHotBoards(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	hbcache := []ptttype.BidInStore{9, 0, 7}
	copy(cache.Shm.Shm.HBcache[:], hbcache)
	cache.Shm.Shm.NHOTs = 3

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
			args:            args{user: testUserecRaw1, uid: 1},
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
			for _, each := range gotSummary {
				each.Total = 0
				each.LastPostTime = 0
			}
			testutil.TDeepEqual(t, "got", gotSummary, tt.expectedSummary)
		})
		wg.Wait()
	}
}

func TestLoadBoardsByBids(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	bids := []ptttype.Bid{10, 1, 8}

	type args struct {
		user *ptttype.UserecRaw
		uid  ptttype.UID
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

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, err := LoadBoardsByBids(tt.args.user, tt.args.uid, tt.args.bids)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardsByBids() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, each := range gotSummaries {
				each.Total = 0
				each.LastPostTime = 0
			}

			testutil.TDeepEqual(t, "got", gotSummaries, tt.expectedSummaries)
		})
		wg.Wait()
	}
}

func TestFindBoardStartIdxByName(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	type args struct {
		boardID *ptttype.BoardID_t
		isAsc   bool
	}
	tests := []struct {
		name             string
		args             args
		expectedStartIdx ptttype.SortIdx
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			args:             args{boardID: testBoardSummary10.Brdname, isAsc: true},
			expectedStartIdx: 12,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotStartIdx, err := FindBoardStartIdxByName(tt.args.boardID, tt.args.isAsc)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindBoardStartIdxByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStartIdx, tt.expectedStartIdx) {
				t.Errorf("FindBoardStartIdxByName() = %v, want %v", gotStartIdx, tt.expectedStartIdx)
			}
		})
		wg.Wait()
	}
}

func TestFindBoardStartIdxByClass(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	type args struct {
		cls     []byte
		boardID *ptttype.BoardID_t
		isAsc   bool
	}
	tests := []struct {
		name             string
		args             args
		expectedStartIdx ptttype.SortIdx
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			args:             args{cls: testBoardSummary10.Title[:4], boardID: testBoardSummary10.Brdname, isAsc: true},
			expectedStartIdx: 12,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotStartIdx, err := FindBoardStartIdxByClass(tt.args.cls, tt.args.boardID, tt.args.isAsc)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindBoardStartIdxByClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStartIdx, tt.expectedStartIdx) {
				t.Errorf("FindBoardStartIdxByClass() = %v, want %v", gotStartIdx, tt.expectedStartIdx)
			}
		})
		wg.Wait()
	}
}

func TestLoadAutoCompleteBoards(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	if cmbbs.Sem.SemID == 0 {
		logrus.Errorf("TestLoadAutoCompleteBoards: before ReloadBCache: semid is invalid")
	}

	cache.ReloadBCache()

	if cmbbs.Sem.SemID == 0 {
		logrus.Errorf("TestLoadAutoCompleteBoards: semid is invalid")
	}

	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.UID
		startIdx ptttype.SortIdx
		nBoards  int
		keyword  []byte
		isAsc    bool
	}
	tests := []struct {
		name                string
		args                args
		expectedSummaries   []*ptttype.BoardSummaryRaw
		expectedNextSummary *ptttype.BoardSummaryRaw
		wantErr             bool
	}{
		// TODO: Add test cases.
		{
			args:                args{user: testUserecRaw4, uid: 11, startIdx: 3, nBoards: 200, keyword: []byte{'a'}, isAsc: true},
			expectedSummaries:   []*ptttype.BoardSummaryRaw{testBoardSummary6},
			expectedNextSummary: nil,
		},
		{
			args:                args{user: testUserecRaw4, uid: 11, startIdx: 4, nBoards: 200, keyword: []byte{'a'}, isAsc: false},
			expectedSummaries:   []*ptttype.BoardSummaryRaw{testBoardSummary6},
			expectedNextSummary: nil,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, gotNextSummary, err := LoadAutoCompleteBoards(tt.args.user, tt.args.uid, tt.args.startIdx, tt.args.nBoards, tt.args.keyword, tt.args.isAsc)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAutoCompleteBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSummaries, tt.expectedSummaries) {
				t.Errorf("LoadAutoCompleteBoards() gotSummaries = %v, want %v", gotSummaries, tt.expectedSummaries)
			}
			if !reflect.DeepEqual(gotNextSummary, tt.expectedNextSummary) {
				t.Errorf("LoadAutoCompleteBoards() gotNextSummary = %v, want %v", gotNextSummary, tt.expectedNextSummary)
			}
		})
		wg.Wait()
	}
}

func TestFindBoardAutoCompleteStartIdx(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	type args struct {
		keyword []byte
		isAsc   bool
	}
	tests := []struct {
		name             string
		args             args
		expectedStartIdx ptttype.SortIdx
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			args:             args{keyword: []byte{'w'}, isAsc: false},
			expectedStartIdx: 12,
		},
		{
			args:             args{keyword: []byte{'y'}, isAsc: false},
			expectedStartIdx: -1,
		},
		{
			args:             args{keyword: []byte{'w'}, isAsc: true},
			expectedStartIdx: 12,
		},
		{
			args:             args{keyword: []byte{'y'}, isAsc: true},
			expectedStartIdx: -1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotStartIdx, err := FindBoardAutoCompleteStartIdx(tt.args.keyword, tt.args.isAsc)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindBoardAutoCompleteStartIdx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStartIdx, tt.expectedStartIdx) {
				t.Errorf("FindBoardAutoCompleteStartIdx() = %v, want %v", gotStartIdx, tt.expectedStartIdx)
			}
		})
		wg.Wait()
	}
}

func TestLoadClassBoards(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.UID
		classBid ptttype.Bid
		bsortBy  ptttype.BSortBy
	}
	tests := []struct {
		name              string
		args              args
		expectedSummaries []*ptttype.BoardSummaryRaw
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			args: args{user: testUserecRaw1, uid: 1, classBid: 1, bsortBy: ptttype.BSORT_BY_NAME},

			expectedSummaries: []*ptttype.BoardSummaryRaw{testClassSummary2, testClassSummary5},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, err := LoadClassBoards(tt.args.user, tt.args.uid, tt.args.classBid, tt.args.bsortBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadClassBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotSummaries, tt.expectedSummaries)
		})
		wg.Wait()
	}
}

func TestLoadFullClassBoards(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.UID
		startBid ptttype.Bid
		nBoards  int
	}
	tests := []struct {
		name                string
		args                args
		expectedSummaries   []*ptttype.BoardSummaryRaw
		expectedNextSummary *ptttype.BoardSummaryRaw
		wantErr             bool
	}{
		// TODO: Add test cases.
		{
			args:              args{user: testUserecRaw1, uid: 1, startBid: 1, nBoards: 100},
			expectedSummaries: []*ptttype.BoardSummaryRaw{testClassSummary2, testClassSummary5},
		},
		{
			args:                args{user: testUserecRaw1, uid: 1, startBid: 1, nBoards: 1},
			expectedSummaries:   []*ptttype.BoardSummaryRaw{testClassSummary2},
			expectedNextSummary: testClassSummary5,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, gotNextSummary, err := LoadFullClassBoards(tt.args.user, tt.args.uid, tt.args.startBid, tt.args.nBoards)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFullClassBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSummaries, tt.expectedSummaries) {
				t.Errorf("LoadFullClassBoards() gotSummaries = %v, want %v", gotSummaries, tt.expectedSummaries)
			}
			if !reflect.DeepEqual(gotNextSummary, tt.expectedNextSummary) {
				t.Errorf("LoadFullClassBoards() gotNextSummary = %v, want %v", gotNextSummary, tt.expectedNextSummary)
			}
		})
		wg.Wait()
	}
}

func TestLoadBoardDetail(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	type args struct {
		user *ptttype.UserecRaw
		uid  ptttype.UID
		bid  ptttype.Bid
	}
	tests := []struct {
		name          string
		args          args
		expectedBoard *ptttype.BoardDetailRaw
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			args:          args{user: testUserecRaw1, uid: 1, bid: 10},
			expectedBoard: testBoardDetailRaw10,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotBoard, err := LoadBoardDetail(tt.args.user, tt.args.uid, tt.args.bid)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotBoard.Total = 0
			gotBoard.LastPostTime = 0
			testutil.TDeepEqual(t, "got", gotBoard, tt.expectedBoard)
		})
		wg.Wait()
	}
}
