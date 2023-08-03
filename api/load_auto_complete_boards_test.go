package api

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadAutoCompleteBoards(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	params := &LoadAutoCompleteBoardsParams{
		StartIdx: "ALLPOST",
		NBoards:  4,
		Keyword:  "a",
		Asc:      true,
	}

	expected := &LoadGeneralBoardsResult{
		Boards:  []*bbs.BoardSummary{testBoardSummary6},
		NextIdx: "",
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
			args:     args{remoteAddr: testIP, uuserID: "SYSOP", params: params},
			expected: expected,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := LoadAutoCompleteBoards(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAutoCompleteBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
		wg.Wait()
	}
}
