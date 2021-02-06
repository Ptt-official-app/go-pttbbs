package api

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestLoadAutoCompleteBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadAutoCompleteBoards(tt.args.remoteAddr, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAutoCompleteBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("LoadAutoCompleteBoards() = %v, want %v", got, tt.expected)
			}
		})
	}
}
