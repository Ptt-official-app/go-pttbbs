package api

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadGeneralArticles(t *testing.T) {
	setupTest()
	defer teardownTest()

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
			args: args{
				uuserID: "SYSOP",
				params: &LoadGeneralArticlesParams{
					StartIdx:  "",
					NArticles: 1,
					Desc:      true,
				},
				path: path,
			},
			expectedResult: &LoadGeneralArticlesResult{
				Articles:       []*bbs.ArticleSummary{testArticleSummary1},
				IsNewest:       true,
				NextIdx:        "1607202239@1Vo_M_CDSYSOP",
				NextCreateTime: 1607202239,
				StartNumIdx:    2,
			},
		},
		{
			args: args{
				uuserID: "SYSOP",
				params: &LoadGeneralArticlesParams{
					StartIdx:  "1607203395@1Vo_f30DSYSOP",
					NArticles: 2,
					Desc:      true,
				},
				path: path,
			},
			expectedResult: &LoadGeneralArticlesResult{
				Articles:    []*bbs.ArticleSummary{testArticleSummary1, testArticleSummary0},
				IsNewest:    true,
				NextIdx:     "",
				StartNumIdx: 2,
			},
		},
		{
			name: "start-idx: 2, n-articles: 1",
			args: args{
				uuserID: "SYSOP",
				params: &LoadGeneralArticlesParams{
					StartIdx:  "1607203395@1Vo_f30DSYSOP",
					NArticles: 1,
					Desc:      true,
				},
				path: path,
			},
			expectedResult: &LoadGeneralArticlesResult{
				Articles:       []*bbs.ArticleSummary{testArticleSummary1},
				IsNewest:       true,
				NextIdx:        "1607202239@1Vo_M_CDSYSOP",
				NextCreateTime: 1607202239,
				StartNumIdx:    2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := LoadGeneralArticles(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotResult, tt.expectedResult)
		})
	}
}
