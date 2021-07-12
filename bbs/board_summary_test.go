package bbs

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestNewBoardSummaryFromRaw(t *testing.T) {
	setupTest()
	defer teardownTest()

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
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)

		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := NewBoardSummaryFromRaw(tt.args.boardSummaryRaw); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("NewBoardSummaryFromRaw() = %v, want %v", got, tt.expected)
			}
		})
		wg.Wait()
	}
}
