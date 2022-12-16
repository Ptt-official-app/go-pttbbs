package cache

import (
	"encoding/binary"
	"io"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func TestGetBCache(t *testing.T) {
	setupTest()
	defer teardownTest()

	boards := []ptttype.BoardHeaderRaw{testBoardHeader0, testBoardHeader1, testBoardHeader2}
	copy(Shm.Shm.BCache[:], boards[:])

	type args struct {
		bidInCache ptttype.Bid
	}
	tests := []struct {
		name          string
		args          args
		expectedBoard *ptttype.BoardHeaderRaw
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			args:          args{1},
			expectedBoard: &testBoardHeader0,
		},
		{
			args:          args{2},
			expectedBoard: &testBoardHeader1,
		},
		{
			args:          args{3},
			expectedBoard: &testBoardHeader2,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotBoard, err := GetBCache(tt.args.bidInCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBoard, tt.expectedBoard) {
				t.Errorf("GetBCache() = %v, want %v", gotBoard, tt.expectedBoard)
			}
		})
		wg.Wait()
	}
}

func TestIsHiddenBoardFriend(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = LoadUHash()

	ReloadBCache()

	type args struct {
		bidInCache ptttype.BidInStore
		uidInCache ptttype.UIDInStore
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{0, 0}, // board: SYSOP user: SYSOP
			expected: true,
		},
		{
			args:     args{0, 1}, // board: SYSOP user: CodingMan
			expected: false,
		},
		{
			args:     args{0, 2}, // board: SYSOP user: pichu
			expected: false,
		},
		{
			args:     args{0, 3}, // board: SYSOP user: Kahou
			expected: true,
		},
		{
			args:     args{0, 4}, // board: SYSOP user: Kahou2
			expected: false,
		},
		{
			args:     args{0, 5}, // board: SYSOP user: (non-exist)
			expected: false,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := IsHiddenBoardFriend(tt.args.bidInCache, tt.args.uidInCache); got != tt.expected {
				t.Errorf("IsHiddenBoardFriend() = %v, want %v", got, tt.expected)
			}
		})
		wg.Wait()
	}
}

func TestNumBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	tests := []struct {
		name     string
		expected int32
	}{
		// TODO: Add test cases.
		{
			expected: 12,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := NumBoards(); got != tt.expected {
				t.Errorf("NumBoards() = %v, want %v", got, tt.expected)
			}
		})
		wg.Wait()
	}
}

func TestReloadBCache(t *testing.T) {
	setupTest()
	defer teardownTest()

	tests := []struct {
		name                  string
		expectedNBoard        int32
		expectedSortedByName  []ptttype.BidInStore
		expectedSortedByClass []ptttype.BidInStore
		expectedBCacheName    []ptttype.BoardID_t
		expectedBCacheTitle   []ptttype.BoardTitle_t
	}{
		// TODO: Add test cases.
		{
			expectedNBoard:        int32(12),
			expectedSortedByName:  testSortedByName,
			expectedSortedByClass: testSortedByClass,
			expectedBCacheName:    testBCacheName,
			expectedBCacheTitle:   testBCacheTitle,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			ReloadBCache()

			nBoard := Shm.Shm.BNumber

			if !reflect.DeepEqual(nBoard, tt.expectedNBoard) {
				t.Errorf("ReloadBCache() = %v, want %v", nBoard, tt.expectedNBoard)
			}

			bsorted := &Shm.Shm.BSorted

			for idx := int32(0); idx < nBoard; idx++ {
				board, _ := GetBCache(ptttype.Bid(idx + 1))
				if types.Cstrcmp(board.Brdname[:], tt.expectedBCacheName[idx][:]) != 0 {
					t.Errorf("bcacheName (%v/%v) = %v, want %v", idx, nBoard, string(board.Brdname[:]), string(tt.expectedBCacheName[idx][:]))
				}
				if types.Cstrcmp(board.Title[:], tt.expectedBCacheTitle[idx][:]) != 0 {
					t.Errorf("bcacheTitle (%v/%v) = %v, want %v", idx, nBoard, string(board.Title[:]), string(tt.expectedBCacheTitle[idx][:]))
				}
			}

			bsortedByName := bsorted[ptttype.BSORT_BY_NAME][:nBoard]
			bsortedByClass := bsorted[ptttype.BSORT_BY_CLASS][:nBoard]
			if !reflect.DeepEqual(bsortedByName, tt.expectedSortedByName) {
				t.Errorf("bsorted-by-name = %v, want %v", bsortedByName, tt.expectedSortedByName)
			}
			if !reflect.DeepEqual(bsortedByClass, tt.expectedSortedByClass) {
				t.Errorf("bsorted-by-class = %v, want %v", bsortedByClass, tt.expectedSortedByClass)
			}
		})
		wg.Wait()
	}
}

func Test_reloadCacheLoadBottom(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	tests := []struct {
		name     string
		expected uint8
	}{
		// TODO: Add test cases.
		{
			expected: 1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			reloadCacheLoadBottom()
			nBottom := Shm.Shm.NBottom[9]

			if nBottom != tt.expected {
				t.Errorf("nBottom: %v want: %v", nBottom, tt.expected)
			}
		})
		wg.Wait()
	}
}

func TestGetBTotal(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	bid1 := ptttype.Bid(3)
	bid1InCache := bid1.ToBidInStore()
	total1 := int32(5)

	Shm.Shm.Total[bid1InCache] = total1

	type args struct {
		bid ptttype.Bid
	}
	tests := []struct {
		name          string
		args          args
		expectedTotal int32
	}{
		// TODO: Add test cases.
		{
			args:          args{bid1},
			expectedTotal: total1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if gotTotal := GetBTotal(tt.args.bid); gotTotal != tt.expectedTotal {
				t.Errorf("GetBTotal() = %v, want %v", gotTotal, tt.expectedTotal)
			}
		})
		wg.Wait()
	}
}

func TestSetBTotal(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	type args struct {
		bid ptttype.Bid
	}
	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		expectedTotal        int32
		expectedLastPostTime types.Time4
	}{
		// TODO: Add test cases.
		{
			args:                 args{10},
			expectedTotal:        2,
			expectedLastPostTime: 1607203395,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := SetBTotal(tt.args.bid); (err != nil) != tt.wantErr {
				t.Errorf("SetBTotal() error = %v, wantErr %v", err, tt.wantErr)
			}

			total := GetBTotal(tt.args.bid)
			if total != tt.expectedTotal {
				t.Errorf("SetBTotal: total: %v want: %v", total, tt.expectedTotal)
			}

			bidInCache := tt.args.bid.ToBidInStore()
			lastPostTime := Shm.Shm.LastPostTime[bidInCache]
			if lastPostTime != tt.expectedLastPostTime {
				t.Errorf("SetBTotal: lastPostTime: %v want: %v", lastPostTime, tt.expectedLastPostTime)
			}
		})
		wg.Wait()
	}
}

func TestSetBottomTotal(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	type args struct {
		bid ptttype.Bid
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		expectedTotal uint8
	}{
		// TODO: Add test cases.
		{
			args:          args{10},
			expectedTotal: 1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := SetBottomTotal(tt.args.bid); (err != nil) != tt.wantErr {
				t.Errorf("SetBottomTotal() error = %v, wantErr %v", err, tt.wantErr)
			}

			bidInCache := tt.args.bid.ToBidInStore()
			total := Shm.Shm.NBottom[bidInCache]
			if total != tt.expectedTotal {
				t.Errorf("SetBottomTotal: total: %v want: %v", total, tt.expectedTotal)
			}
		})
		wg.Wait()
	}
}

func TestGetBid(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	boardID0 := &ptttype.BoardID_t{}
	copy(boardID0[:], []byte("WhoAmI"))

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("SYSOP"))

	type args struct {
		boardID *ptttype.BoardID_t
	}
	tests := []struct {
		name        string
		args        args
		expectedBid ptttype.Bid
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{boardID: boardID0},
			expectedBid: 10,
		},
		{
			args:        args{boardID: boardID1},
			expectedBid: 1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotBid, err := GetBid(tt.args.boardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBid, tt.expectedBid) {
				t.Errorf("GetBid() = %v, want %v", gotBid, tt.expectedBid)
			}
		})
		wg.Wait()
	}
}

func TestFindBoardIdxByName(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	boardID0 := &ptttype.BoardID_t{}
	copy(boardID0[:], []byte("WhoAmI"))

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("SYSOP"))

	boardID2 := &ptttype.BoardID_t{}
	copy(boardID2[:], []byte("2..........."))

	boardID3 := &ptttype.BoardID_t{}
	copy(boardID3[:], []byte("EditExp"))

	boardID4 := &ptttype.BoardID_t{}
	copy(boardID4[:], []byte("WhoAmA"))

	boardID5 := &ptttype.BoardID_t{}
	copy(boardID5[:], []byte("junk"))

	boardID6 := &ptttype.BoardID_t{}
	copy(boardID6[:], []byte("Note"))

	boardID7 := &ptttype.BoardID_t{}
	copy(boardID7[:], []byte("Record"))

	boardID8 := &ptttype.BoardID_t{}
	copy(boardID8[:], []byte("WhoAmJ"))

	boardID9 := &ptttype.BoardID_t{}
	copy(boardID9[:], []byte("a0"))

	type args struct {
		boardID *ptttype.BoardID_t
		isAsc   bool
	}
	tests := []struct {
		name        string
		args        args
		expectedIdx ptttype.SortIdx
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{boardID: boardID0, isAsc: true},
			expectedIdx: 12,
		},
		{
			args:        args{boardID: boardID1, isAsc: true},
			expectedIdx: 11,
		},
		{
			args:        args{boardID: boardID2, isAsc: true},
			expectedIdx: 2,
		},
		{
			args:        args{boardID: boardID3, isAsc: true},
			expectedIdx: 6,
		},
		{
			args:        args{boardID: boardID4, isAsc: true},
			expectedIdx: 12,
		},
		{
			args:        args{boardID: boardID5, isAsc: true},
			expectedIdx: 7,
		},
		{
			args:        args{boardID: boardID6, isAsc: true},
			expectedIdx: 8,
		},
		{
			args:        args{boardID: boardID7, isAsc: true},
			expectedIdx: 9,
		},
		{
			args:        args{boardID: boardID8, isAsc: true},
			expectedIdx: -1,
		},
		{
			args:        args{boardID: boardID9, isAsc: true},
			expectedIdx: 3,
		},

		{
			args:        args{boardID: boardID0, isAsc: false},
			expectedIdx: 12,
		},
		{
			args:        args{boardID: boardID1, isAsc: false},
			expectedIdx: 11,
		},
		{
			args:        args{boardID: boardID2, isAsc: false},
			expectedIdx: 2,
		},
		{
			args:        args{boardID: boardID3, isAsc: false},
			expectedIdx: 6,
		},
		{
			args:        args{boardID: boardID4, isAsc: false},
			expectedIdx: 11,
		},
		{
			args:        args{boardID: boardID5, isAsc: false},
			expectedIdx: 7,
		},
		{
			args:        args{boardID: boardID6, isAsc: false},
			expectedIdx: 8,
		},
		{
			args:        args{boardID: boardID7, isAsc: false},
			expectedIdx: 9,
		},
		{
			args:        args{boardID: boardID8, isAsc: false},
			expectedIdx: 12,
		},
		{
			args:        args{boardID: boardID9, isAsc: false},
			expectedIdx: 2,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotIdx, err := FindBoardIdxByName(tt.args.boardID, tt.args.isAsc)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindBoardIdxByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIdx, tt.expectedIdx) {
				t.Errorf("FindBoardIdxByName() = %v, want %v", gotIdx, tt.expectedIdx)
			}
		})
		wg.Wait()
	}
}

func TestFindBoardIdxByClass(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	testTitle0 := make([]byte, 4)
	copy(testTitle0, testBCacheTitle[0][:])

	testBrdname0 := &ptttype.BoardID_t{}
	copy(testBrdname0[:], []byte("WhoAMA"))

	testTitle1 := make([]byte, 4)
	copy(testTitle1, testBCacheTitle[0][:])
	testTitle1[3]--

	testTitle2 := make([]byte, 4)
	copy(testTitle2, testBCacheTitle[0][:])
	testTitle2[3]++

	type args struct {
		cls     []byte
		boardID *ptttype.BoardID_t
		isAsc   bool
	}
	tests := []struct {
		name        string
		args        args
		expectedIdx ptttype.SortIdx
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{cls: testBCacheTitle[0][:4], boardID: &testBCacheName[0], isAsc: true},
			expectedIdx: 11,
		},
		{
			args:        args{cls: testBCacheTitle[1][:4], boardID: &testBCacheName[1], isAsc: true},
			expectedIdx: 1,
		},
		{
			args:        args{cls: testBCacheTitle[2][:4], boardID: &testBCacheName[2], isAsc: true},
			expectedIdx: 3,
		},
		{
			args:        args{cls: testBCacheTitle[3][:4], boardID: &testBCacheName[3], isAsc: true},
			expectedIdx: 4,
		},
		{
			args:        args{cls: testBCacheTitle[4][:4], boardID: &testBCacheName[4], isAsc: true},
			expectedIdx: 2,
		},
		{
			args:        args{cls: testBCacheTitle[5][:4], boardID: &testBCacheName[5], isAsc: true},
			expectedIdx: 6,
		},
		{
			args:        args{cls: testBCacheTitle[6][:4], boardID: &testBCacheName[6], isAsc: true},
			expectedIdx: 7,
		},
		{
			args:        args{cls: testBCacheTitle[7][:4], boardID: &testBCacheName[7], isAsc: true},
			expectedIdx: 9,
		},
		{
			args:        args{cls: testBCacheTitle[8][:4], boardID: &testBCacheName[8], isAsc: true},
			expectedIdx: 10,
		},
		{
			args:        args{cls: testBCacheTitle[9][:4], boardID: &testBCacheName[9], isAsc: true},
			expectedIdx: 12,
		},
		{
			args:        args{cls: testBCacheTitle[10][:4], boardID: &testBCacheName[10], isAsc: true},
			expectedIdx: 8,
		},
		{
			args:        args{cls: testBCacheTitle[11][:4], boardID: &testBCacheName[11], isAsc: true},
			expectedIdx: 5,
		},
		{
			name:        "title0-brdname0, WhoAMA, asc",
			args:        args{cls: testTitle0, boardID: testBrdname0, isAsc: true},
			expectedIdx: 12,
		},
		{
			name:        "title1-brdname0, between Security/AllHIDPOST, asc",
			args:        args{cls: testTitle1, boardID: testBrdname0, isAsc: true},
			expectedIdx: 5,
		},
		{
			name:        "title2-brdname0, after WhoAmI, asc",
			args:        args{cls: testTitle2, boardID: testBrdname0, isAsc: true},
			expectedIdx: -1,
		},

		{
			args:        args{cls: testBCacheTitle[0][:4], boardID: &testBCacheName[0], isAsc: false},
			expectedIdx: 11,
		},
		{
			args:        args{cls: testBCacheTitle[1][:4], boardID: &testBCacheName[1], isAsc: false},
			expectedIdx: 1,
		},
		{
			args:        args{cls: testBCacheTitle[2][:4], boardID: &testBCacheName[2], isAsc: false},
			expectedIdx: 3,
		},
		{
			args:        args{cls: testBCacheTitle[3][:4], boardID: &testBCacheName[3], isAsc: false},
			expectedIdx: 4,
		},
		{
			args:        args{cls: testBCacheTitle[4][:4], boardID: &testBCacheName[4], isAsc: false},
			expectedIdx: 2,
		},
		{
			args:        args{cls: testBCacheTitle[5][:4], boardID: &testBCacheName[5], isAsc: false},
			expectedIdx: 6,
		},
		{
			args:        args{cls: testBCacheTitle[6][:4], boardID: &testBCacheName[6], isAsc: false},
			expectedIdx: 7,
		},
		{
			args:        args{cls: testBCacheTitle[7][:4], boardID: &testBCacheName[7], isAsc: false},
			expectedIdx: 9,
		},
		{
			args:        args{cls: testBCacheTitle[8][:4], boardID: &testBCacheName[8], isAsc: false},
			expectedIdx: 10,
		},
		{
			args:        args{cls: testBCacheTitle[9][:4], boardID: &testBCacheName[9], isAsc: false},
			expectedIdx: 12,
		},
		{
			args:        args{cls: testBCacheTitle[10][:4], boardID: &testBCacheName[10], isAsc: false},
			expectedIdx: 8,
		},
		{
			args:        args{cls: testBCacheTitle[11][:4], boardID: &testBCacheName[11], isAsc: false},
			expectedIdx: 5,
		},
		{
			name:        "title0-brdname0, WhoAMA, desc",
			args:        args{cls: testTitle0, boardID: testBrdname0, isAsc: false},
			expectedIdx: 11,
		},
		{
			name:        "title1-brdname0, between Security/AllHIDPOST, desc",
			args:        args{cls: testTitle1, boardID: testBrdname0, isAsc: false},
			expectedIdx: 4,
		},
		{
			name:        "title2-brdname0, after WhoAmI, desc",
			args:        args{cls: testTitle2, boardID: testBrdname0, isAsc: false},
			expectedIdx: 12,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotIdx, err := FindBoardIdxByClass(tt.args.cls, tt.args.boardID, tt.args.isAsc)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindBoardIdxByClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIdx, tt.expectedIdx) {
				t.Errorf("FindBoardIdxByClass() = %v, want %v", gotIdx, tt.expectedIdx)
			}
		})
		wg.Wait()
	}
}

func TestFindBoardAutoCompleteStartIdx(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

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
			name:             "a, asc (ALLHIDPOST)",
			args:             args{keyword: []byte{'a'}, isAsc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "a, desc (ALLPOST)",
			args:             args{keyword: []byte{'a'}, isAsc: false},
			expectedStartIdx: 4,
		},
		{
			args:             args{keyword: []byte{'b'}, isAsc: true},
			expectedStartIdx: -1,
		},
		{
			args:             args{keyword: []byte{'r'}, isAsc: true},
			expectedStartIdx: 9,
		},
		{
			args:             args{keyword: []byte{'r'}, isAsc: false},
			expectedStartIdx: 9,
		},
		{
			name:             "s, asc (Security)",
			args:             args{keyword: []byte{'s'}, isAsc: true},
			expectedStartIdx: 10,
		},
		{
			name:             "s, desc (SYSOP)",
			args:             args{keyword: []byte{'s'}, isAsc: false},
			expectedStartIdx: 11,
		},
		{
			args:             args{keyword: []byte{'t'}, isAsc: true},
			expectedStartIdx: -1,
		},
		{
			args:             args{keyword: []byte{'w'}, isAsc: true},
			expectedStartIdx: 12,
		},
		{
			args:             args{keyword: []byte{'w'}, isAsc: false},
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

func TestGetBTotalWithRetry(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	type args struct {
		bid ptttype.Bid
	}
	tests := []struct {
		name          string
		args          args
		expectedTotal int32
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			args:          args{10},
			expectedTotal: 2,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotTotal, err := GetBTotalWithRetry(tt.args.bid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBTotalWithRetry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotal != tt.expectedTotal {
				t.Errorf("GetBTotalWithRetry() = %v, want %v", gotTotal, tt.expectedTotal)
			}
		})
		wg.Wait()
	}
}

func TestNHots(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	tests := []struct {
		name          string
		expectedNhots uint8
	}{
		// TODO: Add test cases.
		{
			expectedNhots: 0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if gotNhots := NHots(); gotNhots != tt.expectedNhots {
				t.Errorf("NHots() = %v, want %v", gotNhots, tt.expectedNhots)
			}
		})
		wg.Wait()
	}
}

func TestSanitizeBMs(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = LoadUHash()

	bms0 := &ptttype.BM_t{}
	copy(bms0[:], []byte("SYSOP/SYSOP2/SYSOP3"))

	expected0 := &ptttype.BM_t{}
	copy(expected0[:], []byte("SYSOP"))

	bms1 := &ptttype.BM_t{}
	copy(bms1[:], []byte("SYSOP/SYSOP2/CodingMan"))

	expected1 := &ptttype.BM_t{}
	copy(expected1[:], []byte("SYSOP/CodingMan"))

	type args struct {
		bms *ptttype.BM_t
	}
	tests := []struct {
		name              string
		args              args
		expectedParsedBMs *ptttype.BM_t
	}{
		// TODO: Add test cases.
		{
			args:              args{bms: bms0},
			expectedParsedBMs: expected0,
		},
		{
			args:              args{bms: bms1},
			expectedParsedBMs: expected1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if gotParsedBMs := SanitizeBMs(tt.args.bms); !reflect.DeepEqual(gotParsedBMs, tt.expectedParsedBMs) {
				t.Errorf("SanitizeBMs() = %v, want %v", gotParsedBMs, tt.expectedParsedBMs)
			}
		})
		wg.Wait()
	}
}

func TestParseBMList(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = LoadUHash()

	bms0 := &ptttype.BM_t{}
	copy(bms0[:], []byte("SYSOP/SYSOP2/SYSOP3"))

	expected0 := &[ptttype.MAX_BMs]ptttype.UID{1, -1, -1, -1}

	bms1 := &ptttype.BM_t{}
	copy(bms1[:], []byte("SYSOP/SYSOP2/CodingMan"))

	expected1 := &[ptttype.MAX_BMs]ptttype.UID{1, 2, -1, -1}

	type args struct {
		bms *ptttype.BM_t
	}
	tests := []struct {
		name         string
		args         args
		expectedUids *[ptttype.MAX_BMs]ptttype.UID
	}{
		// TODO: Add test cases.
		{
			args:         args{bms: bms0},
			expectedUids: expected0,
		},
		{
			args:         args{bms: bms1},
			expectedUids: expected1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if gotUids := ParseBMList(tt.args.bms); !reflect.DeepEqual(gotUids, tt.expectedUids) {
				t.Errorf("ParseBMList() = %v, want %v", gotUids, tt.expectedUids)
			}
		})
		wg.Wait()
	}
}

func TestResetBoard(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = LoadUHash()

	bid0 := ptttype.Bid(1)
	expected0 := &testBoardHeader3

	type args struct {
		bid ptttype.Bid
	}
	tests := []struct {
		name     string
		args     args
		expected *ptttype.BoardHeaderRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{bid: bid0},
			expected: expected0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := ResetBoard(tt.args.bid); (err != nil) != tt.wantErr {
				t.Errorf("ResetBoard() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, _ := GetBCache(tt.args.bid)
			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
		wg.Wait()
	}
}

func Test_buildBMCache(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = LoadUHash()

	Shm.Shm.BCache[0] = testBoardHeader4

	expected0 := []ptttype.UID{1, 2, -1, -1}

	type args struct {
		bid ptttype.Bid
	}
	tests := []struct {
		name     string
		args     args
		expected []ptttype.UID
	}{
		// TODO: Add test cases.
		{
			args:     args{bid: 1},
			expected: expected0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			buildBMCache(tt.args.bid)

			bidInStore := tt.args.bid.ToBidInStore()
			got := &Shm.Shm.BMCache[bidInStore]

			testutil.TDeepEqual(t, "got", got[:], tt.expected)
		})
		wg.Wait()
	}
}

func TestAddbrdTouchCache(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	file, _ := os.OpenFile(ptttype.FN_BOARD, os.O_WRONLY, 0o644)
	defer file.Close()

	_, _ = file.Seek(0, io.SeekEnd)
	_ = types.BinaryWrite(file, binary.LittleEndian, &testBoardHeader13)

	tests := []struct {
		name        string
		expectedBid ptttype.Bid
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			expectedBid: 13,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotBid, err := AddbrdTouchCache()
			if (err != nil) != tt.wantErr {
				t.Errorf("AddbrdTouchCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBid, tt.expectedBid) {
				t.Errorf("AddbrdTouchCache() = %v, want %v", gotBid, tt.expectedBid)
			}
		})
		wg.Wait()
	}
}

func TestSetLastPosttime(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	type args struct {
		bid   ptttype.Bid
		nowTS types.Time4
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected types.Time4
	}{
		// TODO: Add test cases.
		{
			args:     args{bid: 10, nowTS: 1234567890},
			expected: 1234567890,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := SetLastPosttime(tt.args.bid, tt.args.nowTS); (err != nil) != tt.wantErr {
				t.Errorf("SetLastPosttime() error = %v, wantErr %v", err, tt.wantErr)
			}

			lastposttime, _ := GetLastPosttime(tt.args.bid)
			testutil.TDeepEqual(t, "lastposttime", lastposttime, tt.expected)
		})

		wg.Wait()
	}
}

func TestTouchBPostNum(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	type args struct {
		bid   ptttype.Bid
		delta int32
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected int32
	}{
		// TODO: Add test cases.
		{
			args:     args{bid: 10, delta: 1},
			expected: 1,
		},
		{
			args:     args{bid: 10, delta: 2},
			expected: 3,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := TouchBPostNum(tt.args.bid, tt.args.delta); (err != nil) != tt.wantErr {
				t.Errorf("TouchBPostNum() error = %v, wantErr %v", err, tt.wantErr)
			}

			total := GetBTotal(tt.args.bid)
			testutil.TDeepEqual(t, "total", total, tt.expected)
		})
		wg.Wait()
	}
}

func TestResolveBoardGroup(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	board, _ := GetBCache(1)
	logrus.Infof("before test: cls: 1: board: %v", board)

	type args struct {
		gid     ptttype.Bid
		bsortBy ptttype.BSortBy
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool

		expectedBid        ptttype.Bid
		expectedFirstChild ptttype.Bid
		expectedNext       ptttype.Bid
		expectedParent     ptttype.Bid
	}{
		// TODO: Add test cases.
		{
			args:               args{gid: 1, bsortBy: 0},
			expectedBid:        1,
			expectedFirstChild: 2,
		},
		{
			args:           args{gid: 1, bsortBy: 0},
			expectedBid:    2,
			expectedParent: 1,
			expectedNext:   5,
		},
		{
			args:           args{gid: 1, bsortBy: 0},
			expectedBid:    5,
			expectedParent: 1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := ResolveBoardGroup(tt.args.gid, tt.args.bsortBy); (err != nil) != tt.wantErr {
				t.Errorf("ResolveBoardGroup() error = %v, wantErr %v", err, tt.wantErr)
			}

			board, _ := GetBCache(tt.expectedBid)
			testutil.TDeepEqual(t, "first", board.FirstChild[tt.args.bsortBy], tt.expectedFirstChild)
			testutil.TDeepEqual(t, "next", board.Next[tt.args.bsortBy], tt.expectedNext)
		})
		wg.Wait()
	}
}

func TestSetBoardChildCount(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()
	ResolveBoardGroup(1, ptttype.BSORT_BY_NAME)

	board, _ := GetBCache(1)
	logrus.Infof("before test: cls: 1: board: %v", board)

	type args struct {
		bid   ptttype.Bid
		count int32
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		expectedCount int32
	}{
		// TODO: Add test cases.
		{
			args:          args{bid: 1, count: 0},
			expectedCount: 0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := SetBoardChildCount(tt.args.bid, tt.args.count); (err != nil) != tt.wantErr {
				t.Errorf("SetBoardChildCount() error = %v, wantErr %v", err, tt.wantErr)
			}
			board, _ := GetBCache(1)

			testutil.TDeepEqual(t, "count", board.ChildCount, tt.expectedCount)
		})
		wg.Wait()
	}
}
