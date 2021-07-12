package ptttype

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestNewArticleSummaryRaw(t *testing.T) {
	aid0 := SortIdx(2)
	boardID0 := &BoardID_t{}
	type args struct {
		aid     SortIdx
		boardID *BoardID_t
		header  *FileHeaderRaw
	}
	tests := []struct {
		name     string
		args     args
		expected *ArticleSummaryRaw
	}{
		// TODO: Add test cases.
		{
			args:     args{aid: aid0, boardID: boardID0, header: testFileHeaderRaw1},
			expected: testArticleSummary1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got := NewArticleSummaryRaw(tt.args.aid, tt.args.boardID, tt.args.header)
			got.BoardID = tt.expected.BoardID

			testutil.TDeepEqual(t, "summary", got, tt.expected)
		})
	}
	wg.Wait()
}
