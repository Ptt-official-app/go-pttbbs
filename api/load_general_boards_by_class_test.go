package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestLoadGeneralBoardsByClass(t *testing.T) {
	setupTest()
	defer teardownTest()

	params := &LoadGeneralBoardsParams{
		StartIdx: "vFSt-Q@ALLPOST",
		NBoards:  4,
		Asc:      true,
	}

	expected := &LoadGeneralBoardsResult{
		Boards:  []*bbs.BoardSummary{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
		NextIdx: "vFSt-Q@Record",
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{uuserID: "SYSOP", params: params},
			expected: expected,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := LoadGeneralBoardsByClass(tt.args.remoteAddr, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoardsByClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("LoadGeneralBoardsByClass() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}
