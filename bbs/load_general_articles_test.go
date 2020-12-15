package bbs

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestLoadGeneralArticles(t *testing.T) {
	//setupTest in for-loop

	type args struct {
		userID      string
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
				userID:      "SYSOP",
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
				userID:      "SYSOP",
				bboardID:    "10_WhoAmI",
				startIdxStr: "2",
				nArticles:   2,
			},
			expectedSummary:    []*ArticleSummary{testArticleSummary0, testArticleSummary1},
			expectedNextIdxStr: "",
			expectedIsNewest:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			gotSummary, gotNextIdxStr, gotIsNewest, err := LoadGeneralArticles(tt.args.userID, tt.args.bboardID, tt.args.startIdxStr, tt.args.nArticles)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for idx, each := range gotSummary {
				if idx >= len(tt.expectedSummary) {
					t.Errorf("LoadGeneralArticles: (%v/%v) %v", idx, len(gotSummary), each)
					continue
				}
				types.TDeepEqual(t, each, tt.expectedSummary[idx])
			}
			if gotNextIdxStr != tt.expectedNextIdxStr {
				t.Errorf("LoadGeneralArticles() gotNextIdxStr = %v, want %v", gotNextIdxStr, tt.expectedNextIdxStr)
			}

			if gotIsNewest != tt.expectedIsNewest {
				t.Errorf("LoadGeneralArticles() isNewest = %v, want %v", gotIsNewest, tt.expectedIsNewest)
			}
		})
	}
}
