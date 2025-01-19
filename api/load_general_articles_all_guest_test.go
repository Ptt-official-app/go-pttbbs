package api

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/gin-gonic/gin"
)

func TestLoadGeneralArticlesAllGuest(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	path := &LoadGeneralArticlesPath{
		BBoardID: "WhoAmI",
	}

	type args struct {
		remoteAddr string
		params     interface{}
		path       interface{}
		c          *gin.Context
	}
	tests := []struct {
		name       string
		args       args
		wantResult interface{}
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				params: &LoadGeneralArticlesParams{
					StartIdx:  "",
					NArticles: 1,
					Desc:      true,
				},
				path: path,
			},
			wantResult: &LoadGeneralArticlesResult{
				Articles:       []*bbs.ArticleSummary{testArticleSummary1AllGuest},
				IsNewest:       true,
				NextIdx:        "1607202239@1Vo_M_CD",
				NextCreateTime: 1607202239,
				StartNumIdx:    2,
			},
		},
		{
			args: args{
				params: &LoadGeneralArticlesParams{
					StartIdx:  "1607203395@1Vo_f30D",
					NArticles: 2,
					Desc:      true,
				},
				path: path,
			},
			wantResult: &LoadGeneralArticlesResult{
				Articles:    []*bbs.ArticleSummary{testArticleSummary1AllGuest, testArticleSummary0AllGuest},
				IsNewest:    true,
				NextIdx:     "",
				StartNumIdx: 2,
			},
		},
		{
			name: "start-idx: 2, n-articles: 1",
			args: args{
				params: &LoadGeneralArticlesParams{
					StartIdx:  "1607203395@1Vo_f30D",
					NArticles: 1,
					Desc:      true,
				},
				path: path,
			},
			wantResult: &LoadGeneralArticlesResult{
				Articles:       []*bbs.ArticleSummary{testArticleSummary1AllGuest},
				IsNewest:       true,
				NextIdx:        "1607202239@1Vo_M_CD",
				NextCreateTime: 1607202239,
				StartNumIdx:    2,
			},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := LoadGeneralArticlesAllGuest(tt.args.remoteAddr, tt.args.params, tt.args.path, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticlesAllGuest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotResult, tt.wantResult)
		})
	}
	wg.Wait()
}
