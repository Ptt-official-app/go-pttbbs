package ptt

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func Test_mNewbrd(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	testBrdname0 := &ptttype.BoardID_t{}
	copy(testBrdname0[:], []byte("mnewboard0"))
	testBrdClass0 := []byte("CPBL")
	testBrdTitle0 := []byte("new-board")

	expectedTitle0 := &ptttype.BoardTitle_t{}
	copy(expectedTitle0[:], []byte("CPBL \xa1\xb7new-board"))
	expected0 := &ptttype.BoardHeaderRaw{
		Brdname: *testBrdname0,
		Title:   *expectedTitle0,
		BrdAttr: 0x200000,
		Gid:     2,
	}

	type args struct {
		user         *ptttype.UserecRaw
		clsBid       ptttype.Bid
		brdname      *ptttype.BoardID_t
		brdClass     []byte
		brdTitle     []byte
		bms          *ptttype.BM_t
		brdAttr      ptttype.BrdAttr
		level        ptttype.PERM
		chessCountry ptttype.ChessCode
		isGroup      bool
		isRecover    bool
	}
	tests := []struct {
		name          string
		args          args
		expectedBoard *ptttype.BoardHeaderRaw
		expectedBid   ptttype.Bid
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			args:          args{user: testUserecRaw1, clsBid: 2, brdname: testBrdname0, brdClass: testBrdClass0, brdTitle: testBrdTitle0},
			expectedBoard: expected0,
			expectedBid:   13,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotBoard, gotBid, err := mNewbrd(tt.args.user, tt.args.clsBid, tt.args.brdname, tt.args.brdClass, tt.args.brdTitle, tt.args.bms, tt.args.brdAttr, tt.args.level, tt.args.chessCountry, tt.args.isGroup, tt.args.isRecover)
			if (err != nil) != tt.wantErr {
				t.Errorf("mNewbrd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBoard, tt.expectedBoard) {
				t.Errorf("mNewbrd() gotBoard = %v, want %v", gotBoard, tt.expectedBoard)
			}
			if !reflect.DeepEqual(gotBid, tt.expectedBid) {
				t.Errorf("mNewbrd() gotBid = %v, want %v", gotBid, tt.expectedBid)
			}
		})
		wg.Wait()
	}
}

func Test_addBoardRecord(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	testBoardHeaderRaw0 := &ptttype.BoardHeaderRaw{
		Brdname: ptttype.BoardID_t{'t', 'e', 's', 't', '1', '0', '0'},
		Title: ptttype.BoardTitle_t{
			0xbc, 0x54, 0xad, 0xf9, 0x20, 0xa1, 0xb7, 0xb0, 0xca, 0xba,
			0x41, 0xac, 0xdd, 0xaa, 0x4f, 0xa4, 0xce, 0xba, 0x71, 0xa6,
			0xb1, 0xa7, 0xeb, 0xbd, 0x5a, 0x00, 0xaf, 0xb8, 0xaa, 0xf8,
			0x20, 0x20, 0xa3, 0xad, 0xa1, 0x49, 0x00, 0x6e,
		},
	}

	type args struct {
		board *ptttype.BoardHeaderRaw
	}
	tests := []struct {
		name        string
		args        args
		expectedBid ptttype.Bid
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{board: testBoardHeaderRaw0},
			expectedBid: 13,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotBid, err := addBoardRecord(tt.args.board)
			if (err != nil) != tt.wantErr {
				t.Errorf("addBoardRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBid, tt.expectedBid) {
				t.Errorf("addBoardRecord() = %v, want %v", gotBid, tt.expectedBid)
			}

			got, _ := cache.GetBCache(gotBid)
			testutil.TDeepEqual(t, "got", got, tt.args.board)
		})
		wg.Wait()
	}
}
