package cmsys

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestGetNumRecords(t *testing.T) {
	type args struct {
		filename string
		size     uintptr
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		// TODO: Add test cases.
		{
			args:     args{filename: "testcase/BOARD", size: ptttype.BOARD_HEADER_RAW_SZ},
			expected: 12,
		},
		{
			args:     args{filename: "testcase/not-exist", size: ptttype.BOARD_HEADER_RAW_SZ},
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNumRecords(tt.args.filename, tt.args.size); got != tt.expected {
				t.Errorf("GetNumRecords() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetRecords(t *testing.T) {
	boardID := &ptttype.BoardID_t{}
	copy(boardID[:], []byte("WhoAmI"))
	filename0 := "./testcase/DIR"
	n0 := 100
	expected0 := []*ptttype.ArticleSummaryRaw{testArticleSummary0}

	filename1 := "./testcase/DIR1"
	startAid1 := ptttype.Aid(1)
	n1 := 100
	expected1 := []*ptttype.ArticleSummaryRaw{testArticleSummary1, testArticleSummary2}

	filename2 := "./testcase/DIR"
	startAid2 := ptttype.Aid(2)
	n2 := 100
	var expected2 []*ptttype.ArticleSummaryRaw

	filename3 := "./testcase/DIR"
	startAid3 := ptttype.Aid(3)
	n3 := 100
	var expected3 []*ptttype.ArticleSummaryRaw

	type args struct {
		boardID  *ptttype.BoardID_t
		filename string
		startAid ptttype.Aid
		n        int
	}
	tests := []struct {
		name     string
		args     args
		expected []*ptttype.ArticleSummaryRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{boardID: boardID, filename: filename0, startAid: 1, n: n0},
			expected: expected0,
		},
		{
			args:     args{boardID: boardID, filename: filename1, startAid: startAid1, n: n1},
			expected: expected1,
		},
		{
			args:     args{boardID: boardID, filename: filename2, startAid: startAid2, n: n2},
			expected: expected2,
		},
		{
			args:     args{boardID: boardID, filename: filename3, startAid: startAid3, n: n3},
			expected: expected3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRecords(tt.args.boardID, tt.args.filename, tt.args.startAid, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for idx, each := range got {
				if idx >= len(tt.expected) {
					t.Errorf("GetRecords: (%v/%v) %v", idx, len(got), each)
					continue
				}
				tt.expected[idx].BoardID = boardID
				testutil.TDeepEqual(t, "fileheader", each.FileHeaderRaw, tt.expected[idx].FileHeaderRaw)
			}
		})
	}
}
