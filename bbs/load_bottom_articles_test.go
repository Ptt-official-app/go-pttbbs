package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadBottomArticles(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID  UUserID
		bboardID BBoardID
	}
	tests := []struct {
		name              string
		args              args
		expectedSummaries []*ArticleSummary
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				uuserID:  "SYSOP",
				bboardID: "10_WhoAmI",
			},
			expectedSummaries: []*ArticleSummary{testBottomSummary1},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, err := LoadBottomArticles(tt.args.uuserID, tt.args.bboardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBottomArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "summaries", gotSummaries, tt.expectedSummaries)
		})
	}
	wg.Wait()
}
