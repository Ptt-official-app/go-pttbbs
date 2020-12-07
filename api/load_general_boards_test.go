package api

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestLoadGeneralBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	params := &LoadGeneralBoardsParams{
		StartIdx: 0,
		NBoards:  4,
		Keyword:  nil,
	}

	expected := &LoadGeneralBoardsResult{
		Boards:  []*bbs.BoardSummary{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
		NextIdx: 8,
	}

	type args struct {
		userID string
		params interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{userID: "SYSOP", params: params},
			expected: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadGeneralBoards(tt.args.userID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("LoadGeneralBoards() = %v, want %v", got, tt.expected)
			}
		})
	}
}
