package api

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
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

	result0 := &LoadHotBoardsResult{
		Boards: []*bbs.BoardSummary{testBoardSummary11, testBoardSummary6, testBoardSummary8},
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult interface{}
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{uuserID: "SYSOP", params: nil},
			expectedResult: result0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := LoadHotBoards(tt.args.remoteAddr, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadHotBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("LoadHotBoards() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
	}
}
