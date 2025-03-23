package api

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/gin-gonic/gin"
)

func TestLoadBottomArticlesAllGuest(t *testing.T) {
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
			args: args{remoteAddr: testIP, path: path},
			wantResult: &LoadGeneralArticlesResult{
				Articles: []*bbs.ArticleSummary{testBottomSummary1AllGuest},
			},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := LoadBottomArticlesAllGuest(tt.args.remoteAddr, tt.args.params, tt.args.path, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBottomArticlesAllGuest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotResult, tt.wantResult)
		})
	}
	wg.Wait()
}
