package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestLoadGeneralArticles(t *testing.T) {
	//setupTest in for-loop
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
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607203395@1Vo_f30DSYSOP",
				nArticles:   1,
				isDesc:      true,
			},
			expectedSummary:        []*ArticleSummary{testArticleSummary1},
			expectedNextIdxStr:     "1607202239@1Vo_M_CDSYSOP",
			expectedNextCreateTime: 1607202239,
			expectedIsNewest:       true,
		},
		{
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607203395@1Vo_f30DSYSOP",
				nArticles:   2,
				isDesc:      false,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary1},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
		},
		{
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CDSYSOP",
				nArticles:   1,
				isDesc:      false,
			},
			expectedSummary:        []*ArticleSummary{testArticleSummary0},
			expectedNextIdxStr:     "1607203395@1Vo_f30DSYSOP",
			expectedNextCreateTime: 1607203395,
			expectedIsNewest:       false,
		},
		{
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CDSYSOP",
				nArticles:   2,
				isDesc:      false,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary0, testArticleSummary1},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
		},
		{
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CESYSOP",
				nArticles:   2,
				isDesc:      false,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary0, testArticleSummary1},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
		},
		{
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CESYSOP",
				nArticles:   2,
				isDesc:      true,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary0},
			expectedNextIdxStr: "",
			expectedIsNewest:   false,
		},
		{
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "1607202240@1Vo_N0CDSYSOP",
				nArticles:   2,
				isDesc:      false,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary1},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
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
