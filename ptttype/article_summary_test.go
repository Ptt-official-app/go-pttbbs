package ptttype

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestNewArticleSummaryRaw(t *testing.T) {
	aid0 := Aid(2)
	boardID0 := &BoardID_t{}
	type args struct {
		aid     Aid
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewArticleSummaryRaw(tt.args.aid, tt.args.boardID, tt.args.header)
			got.BoardID = tt.args.boardID

			types.TDeepEqual(t, "summary", got, tt.expected)
		})
	}
}
