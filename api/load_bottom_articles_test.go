package api

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadBottomArticles(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	path := &LoadGeneralArticlesPath{
		BBoardID: "10_WhoAmI",
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
		path       interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult interface{}
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args: args{remoteAddr: testIP, uuserID: "SYSOP", path: path},
			expectedResult: &LoadGeneralArticlesResult{
				Articles: []*bbs.ArticleSummary{testBottomSummary1},
			},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := LoadBottomArticles(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBottomArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotResult, tt.expectedResult)
		})
	}
	wg.Wait()
}
