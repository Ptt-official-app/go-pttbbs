package bbs

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestNewBoardSummaryFromRaw(t *testing.T) {
	type args struct {
		boardSummaryRaw *ptttype.BoardSummaryRaw
	}
	tests := []struct {
		name     string
		args     args
		expected *BoardSummary
	}{
		// TODO: Add test cases.
		{
			args:     args{testBoardSummaryRaw6},
			expected: testBoardSummary6,
		},
	}
	for _, tt := range tests {
		setupTest()
		defer teardownTest()

		t.Run(tt.name, func(t *testing.T) {
			if got := NewBoardSummaryFromRaw(tt.args.boardSummaryRaw); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("NewBoardSummaryFromRaw() = %v, want %v", got, tt.expected)
			}
		})
	}
}
