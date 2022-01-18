package api

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadHotBoards(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

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

	result0 := &LoadHotBoardsResult{
		Boards: []*bbs.BoardSummary{testBoardSummary10, testBoardSummary1, testBoardSummary8},
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult *LoadHotBoardsResult
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{uuserID: "SYSOP", params: nil},
			expectedResult: result0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := LoadHotBoards(tt.args.remoteAddr, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadHotBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result := gotResult.(*LoadHotBoardsResult)
			for idx, each := range result.Boards {
				if idx >= len(tt.expectedResult.Boards) {
					break
				}

				each.StatAttr = tt.expectedResult.Boards[idx].StatAttr
				each.Total = tt.expectedResult.Boards[idx].Total
			}
			testutil.TDeepEqual(t, "got", result, tt.expectedResult)
		})
		wg.Wait()
	}
}
