package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestLoadGeneralArticles(t *testing.T) {
	// setupTest in for-loop
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID     UUserID
		bboardID    BBoardID
		startIdxStr string
		nArticles   int
		isDesc      bool
	}
	tests := []struct {
		name                   string
		args                   args
		expectedSummary        []*ArticleSummary
		expectedNextIdxStr     string
		expectedNextCreateTime types.Time4
		expectedIsNewest       bool
		wantErr                bool
	}{
		// TODO: Add test cases.
		{
			name: "most basic setup: find 1 article with desc, starting from nil",
			args: args{
				uuserID:   "SYSOP",
				bboardID:  "10_WhoAmI",
				nArticles: 1,
				isDesc:    true,
			},
			expectedSummary:        []*ArticleSummary{testArticleSummary1},
			expectedNextIdxStr:     "1607202239@1Vo_M_CD",
			expectedNextCreateTime: 1607202239,
			expectedIsNewest:       true,
		},
		{
			name: "most basic setup: find 2 articles with desc, starting from nil",
			args: args{
				uuserID:   "SYSOP",
				bboardID:  "10_WhoAmI",
				nArticles: 2,
				isDesc:    true,
			},
			expectedSummary:  []*ArticleSummary{testArticleSummary1, testArticleSummary0},
			expectedIsNewest: true,
		},
		{
			name: "most basic setup: find 1 article with asc (not desc), starting from nil",
			args: args{
				uuserID:   "SYSOP",
				bboardID:  "10_WhoAmI",
				nArticles: 1,
				isDesc:    false,
			},
			expectedSummary:        []*ArticleSummary{testArticleSummary0},
			expectedNextIdxStr:     "1607203395@1Vo_f30D",
			expectedNextCreateTime: 1607203395,
		},
		{
			name: "most basic setup: find 2 articles with asc (not desc), starting from nil",
			args: args{
				uuserID:   "SYSOP",
				bboardID:  "10_WhoAmI",
				nArticles: 2,
				isDesc:    false,
			},
			expectedSummary:  []*ArticleSummary{testArticleSummary0, testArticleSummary1},
			expectedIsNewest: true,
		},
		{
			name: "most basic setup: find 1 article with desc, starting from the last article",
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607203395@1Vo_f30D",
				nArticles:   1,
				isDesc:      true,
			},
			expectedSummary:        []*ArticleSummary{testArticleSummary1},
			expectedNextIdxStr:     "1607202239@1Vo_M_CD",
			expectedNextCreateTime: 1607202239,
			expectedIsNewest:       true,
		},
		{
			name: "find 2 articles with asc (not isDesc), starting from the last article, should return 1 article",
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607203395@1Vo_f30D",
				nArticles:   2,
				isDesc:      false,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary1},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
		},
		{
			name: "find 1 article with asc (not isDesc), starting from the beginning article",
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CD",
				nArticles:   1,
				isDesc:      false,
			},
			expectedSummary:        []*ArticleSummary{testArticleSummary0},
			expectedNextIdxStr:     "1607203395@1Vo_f30D",
			expectedNextCreateTime: 1607203395,
			expectedIsNewest:       false,
		},
		{
			name: "find 2 articles with asc (not desc), starting from the beginning article, should return 2 articles",
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CD",
				nArticles:   2,
				isDesc:      false,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary0, testArticleSummary1},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
		},
		{
			name: "find 2 articles with asc (not desc), starting from same create-time of the beginning article, but with diff filename, should return 1 articles",
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CE",
				nArticles:   2,
				isDesc:      false,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary0, testArticleSummary1},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
		},
		{
			name: "find 2 articles with desc, starting from same create-time of the beginning article, but with diff filename, should return err.",
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CE",
				nArticles:   2,
				isDesc:      true,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary0},
			expectedNextIdxStr: "",
		},
		{
			name: "find 2 articles with asc (not desc), starting after the beginning article, should return 1 article (last article)",
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202240@1Vo_N0CD",
				nArticles:   2,
				isDesc:      false,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary1},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
		},
		{
			name: "find 2 articles with desc, starting after the beginning article, should return 1 article (beginning article)",
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202240@1Vo_N0CD",
				nArticles:   2,
				isDesc:      true,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary0},
			expectedNextIdxStr: "",
			expectedIsNewest:   false,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, gotNextIdxStr, gotNextCreateTime, gotIsNewest, _, err := LoadGeneralArticles(tt.args.uuserID, tt.args.bboardID, tt.args.startIdxStr, tt.args.nArticles, tt.args.isDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)

			if gotNextIdxStr != tt.expectedNextIdxStr {
				t.Errorf("LoadGeneralArticles() gotNextIdxStr = %v, want %v", gotNextIdxStr, tt.expectedNextIdxStr)
			}

			if gotNextCreateTime != tt.expectedNextCreateTime {
				t.Errorf("LoadGeneralArticles() gotNextCreateTime = %v, want %v", gotNextCreateTime, tt.expectedNextCreateTime)
			}

			if gotIsNewest != tt.expectedIsNewest {
				t.Errorf("LoadGeneralArticles() isNewest = %v, want %v", gotIsNewest, tt.expectedIsNewest)
			}
		})
		wg.Wait()
	}
}
