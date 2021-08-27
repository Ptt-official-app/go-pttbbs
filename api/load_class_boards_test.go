package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestLoadClassBoards(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	params0 := &LoadClassBoardsParams{IsSystem: true}
	path0 := &LoadClassBoardsPath{ClsID: 1}
	result0 := &LoadClassBoardsResult{
		Boards: []*bbs.BoardSummary{testClassSummary2, testClassSummary5},
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
		path       interface{}
	}
	tests := []struct {
		name        string
		args        args
		expectedRet *LoadClassBoardsResult
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{remoteAddr: testIP, uuserID: "SYSOP", params: params0, path: path0},
			expectedRet: result0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotRet, err := LoadClassBoards(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadClassBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.expectedRet) {
				t.Errorf("LoadClassBoards() = %v, want %v", gotRet, tt.expectedRet)
			}
		})
		wg.Wait()
	}
}
