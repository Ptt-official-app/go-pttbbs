package bbs

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestNewBoardDetailFromRaw(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		boardHeaderRaw *ptttype.BoardHeaderRaw
		bid            ptttype.Bid
	}
	tests := []struct {
		name     string
		args     args
		expected *BoardDetail
	}{
		// TODO: Add test cases.
		{
			args:     args{testBoardHeader3, 1},
			expected: testBoardDetail3,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)

		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			t.Log(NewBoardDetailFromRaw(tt.args.boardHeaderRaw, tt.args.bid))
			t.Log(tt.expected)
			if got := NewBoardDetailFromRaw(tt.args.boardHeaderRaw, tt.args.bid); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("NewBoardDetailFromRaw() = %v, want %v", got, tt.expected)
			}
		})
		wg.Wait()
	}
}
