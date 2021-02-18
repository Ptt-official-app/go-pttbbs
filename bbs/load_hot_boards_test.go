package bbs

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadHotBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	hbcache := []ptttype.BidInStore{10, 5, 7}
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
			expectedSummary: []*BoardSummary{testBoardSummary11, testBoardSummary6, testBoardSummary8},
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

			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)
		})
	}
	wg.Wait()
}
