package ptt

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestLoadGeneralArticles(t *testing.T) {

	boardID := &ptttype.BoardID_t{}
	copy(boardID[:], []byte("WhoAmI"))
	bid := ptttype.Bid(10)

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("SYSOP"))
	bid1 := ptttype.Bid(1)

	type args struct {
		user       *ptttype.UserecRaw
		uid        ptttype.Uid
		boardIDRaw *ptttype.BoardID_t
		bid        ptttype.Bid
		startIdx   ptttype.SortIdx
		nArticles  int
	}
	tests := []struct {
		name               string
		args               args
		expectedSummaryRaw []*ptttype.ArticleSummaryRaw
		expectedNextIdx    ptttype.SortIdx
		expectedIsNewest   bool
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startIdx:   3,
				nArticles:  1,
			},
			expectedSummaryRaw: []*ptttype.ArticleSummaryRaw{testArticleSummary1},
			expectedNextIdx:    1,
			expectedIsNewest:   true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startIdx:   2,
				nArticles:  100,
			},
			expectedSummaryRaw: []*ptttype.ArticleSummaryRaw{testArticleSummary0, testArticleSummary1},
			expectedNextIdx:    -1,
			expectedIsNewest:   true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startIdx:   2,
				nArticles:  2,
			},
			expectedSummaryRaw: []*ptttype.ArticleSummaryRaw{testArticleSummary0, testArticleSummary1},
			expectedNextIdx:    -1,
			expectedIsNewest:   true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startIdx:   2,
				nArticles:  1,
			},
			expectedSummaryRaw: []*ptttype.ArticleSummaryRaw{testArticleSummary1},
			expectedNextIdx:    1,
			expectedIsNewest:   true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startIdx:   1,
				nArticles:  1,
			},
			expectedSummaryRaw: []*ptttype.ArticleSummaryRaw{testArticleSummary0},
			expectedNextIdx:    -1,
			expectedIsNewest:   false,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID1,
				bid:        bid1,
				startIdx:   2,
				nArticles:  1,
			},
			expectedSummaryRaw: nil,
			expectedNextIdx:    -1,
			expectedIsNewest:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			cache.ReloadBCache()

			gotSummaryRaw, gotNextIdx, gotIsNewest, err := LoadGeneralArticles(tt.args.user, tt.args.uid, tt.args.boardIDRaw, tt.args.bid, tt.args.startIdx, tt.args.nArticles)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotSummaryRaw) != len(tt.expectedSummaryRaw) {
				t.Errorf("LoadGeneralArticles: len(got): %v expected: %v", len(gotSummaryRaw), len(tt.expectedSummaryRaw))
			}
			for idx, each := range gotSummaryRaw {
				if idx >= len(tt.expectedSummaryRaw) {
					t.Errorf("LoadGeneralArticles: (%v/%v) %v", idx, len(gotSummaryRaw), each)
					continue
				}
				types.TDeepEqual(t, each.FileHeaderRaw, tt.expectedSummaryRaw[idx].FileHeaderRaw)
			}
			if !reflect.DeepEqual(gotNextIdx, tt.expectedNextIdx) {
				t.Errorf("LoadGeneralArticles() gotNextIdx = %v, want %v", gotNextIdx, tt.expectedNextIdx)
			}
			if !reflect.DeepEqual(gotIsNewest, tt.expectedIsNewest) {
				t.Errorf("LoadGeneralArticles() isNewest = %v, want %v", gotIsNewest, tt.expectedIsNewest)
			}
		})
	}
}
