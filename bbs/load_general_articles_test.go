package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
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
	}
	tests := []struct {
		name               string
		args               args
		expectedSummary    []*ArticleSummary
		expectedNextIdxStr string
		expectedIsNewest   bool
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "2",
				nArticles:   1,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary1},
			expectedNextIdxStr: "1",
			expectedIsNewest:   true,
		},
		{
			args: args{
				uuserID:     "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "2",
				nArticles:   2,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary1, testArticleSummary0},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, gotNextIdxStr, gotIsNewest, err := LoadGeneralArticles(tt.args.uuserID, tt.args.bboardID, tt.args.startIdxStr, tt.args.nArticles)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)

			if gotNextIdxStr != tt.expectedNextIdxStr {
				t.Errorf("LoadGeneralArticles() gotNextIdxStr = %v, want %v", gotNextIdxStr, tt.expectedNextIdxStr)
			}

			if gotIsNewest != tt.expectedIsNewest {
				t.Errorf("LoadGeneralArticles() isNewest = %v, want %v", gotIsNewest, tt.expectedIsNewest)
			}
		})
		wg.Wait()
	}
}
