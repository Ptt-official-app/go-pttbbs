package api

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestLoadGeneralBoardsByClass(t *testing.T) {
	setupTest()
	defer teardownTest()

	params := &LoadGeneralBoardsParams{
		StartIdx: strconv.Itoa(int(0)),
		NBoards:  4,
	}

	expected := &LoadGeneralBoardsResult{
		Boards:  []*bbs.BoardSummary{testBoardSummary6, testBoardSummary7, testBoardSummary11, testBoardSummary8},
		NextIdx: strconv.Itoa(int(9)),
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
}
