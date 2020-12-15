package api

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestLoadGeneralArticles(t *testing.T) {
	setupTest()
	defer teardownTest()

	path := &LoadGeneralArticlesPath{
		BBoardID: "10_WhoAmI",
	}

	type args struct {
		remoteAddr string
		userID     string
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
				userID: "SYSOP",
				params: &LoadGeneralArticlesParams{
					StartIdx:  "",
					NArticles: 1,
				},
				path: path,
			},
			expectedResult: &LoadGeneralArticlesResult{
				Articles: []*bbs.ArticleSummary{testArticleSummary1},
				IsNewest: true,
				NextIdx:  "1",
			},
		},
		{
			args: args{
				userID: "SYSOP",
				params: &LoadGeneralArticlesParams{
					StartIdx:  "2",
					NArticles: 2,
				},
				path: path,
			},
			expectedResult: &LoadGeneralArticlesResult{
				Articles: []*bbs.ArticleSummary{testArticleSummary0, testArticleSummary1},
				IsNewest: true,
				NextIdx:  "",
			},
		},
		{
			args: args{
				userID: "SYSOP",
				params: &LoadGeneralArticlesParams{
					StartIdx:  "2",
					NArticles: 1,
				},
				path: path,
			},
			expectedResult: &LoadGeneralArticlesResult{
				Articles: []*bbs.ArticleSummary{testArticleSummary1},
				IsNewest: true,
				NextIdx:  "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := LoadGeneralArticles(tt.args.remoteAddr, tt.args.userID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("LoadGeneralArticles() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
	}
}
