package api

import (
	"strconv"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadGeneralBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	params := &LoadGeneralBoardsParams{
		StartIdx: strconv.Itoa(int(0)),
		NBoards:  4,
		Asc:      true,
	}

	expected := &LoadGeneralBoardsResult{
		Boards:  []*bbs.BoardSummary{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
		NextIdx: "Record",
	}

	type args struct {
		uuserID bbs.UUserID
		params  interface{}
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
			got, err := LoadGeneralBoards(testIP, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			theGot, _ := got.(*LoadGeneralBoardsResult)
			theExpected, _ := tt.expected.(*LoadGeneralBoardsResult)

			testutil.TDeepEqual(t, "boards", theGot.Boards, theExpected.Boards)
		})
	}
	wg.Wait()
}
