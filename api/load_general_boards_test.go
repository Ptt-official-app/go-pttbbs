package api

import (
	"strconv"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestLoadGeneralBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	params := &LoadGeneralBoardsParams{
		StartIdx: strconv.Itoa(int(0)),
		NBoards:  4,
	}

	expected := &LoadGeneralBoardsResult{
		Boards:  []*bbs.BoardSummary{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
		NextIdx: strconv.Itoa(int(8)),
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadGeneralBoards(testIP, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			theGot, _ := got.(*LoadGeneralBoardsResult)
			theExpected, _ := tt.expected.(*LoadGeneralBoardsResult)

			types.TDeepEqual(t, "boards", theGot.Boards, theExpected.Boards)
		})
	}
}
