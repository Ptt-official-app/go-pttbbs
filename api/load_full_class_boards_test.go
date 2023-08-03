package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestLoadFullClassBoards(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	params0 := &LoadFullClassBoardsParams{
		StartBid: 1,
		NBoards:  100,
		IsSystem: true,
	}

	result0 := &LoadFullClassBoardsResult{
		Boards: []*bbs.BoardSummary{testClassSummary2, testClassSummary5},
	}
	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
	}
	tests := []struct {
		name        string
		args        args
		expectedRet *LoadFullClassBoardsResult
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{remoteAddr: testIP, uuserID: "SYSOP", params: params0},
			expectedRet: result0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotRet, err := LoadFullClassBoards(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFullClassBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.expectedRet) {
				t.Errorf("LoadFullClassBoards() = %v, want %v", gotRet, tt.expectedRet)
			}
		})
		wg.Wait()
	}
}
