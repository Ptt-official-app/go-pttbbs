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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			gotSummary, gotNextIdxStr, err := LoadGeneralArticles(tt.args.userID, tt.args.bboardID, tt.args.startIdxStr, tt.args.nArticles)
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
		})
	}
}
