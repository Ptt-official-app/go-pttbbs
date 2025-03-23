package api

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

func TestGetArticleAllGuest(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	// keep chtime for testing
	filename := "testcase/boards/W/WhoAmI/M.1607202239.A.30D"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	goodParams := &GetArticleParams{}
	goodPath := &GetArticlePath{
		BBoardID:  bbs.BBoardID("WhoAmI"),
		ArticleID: bbs.ArticleID("1Vo_M_CD"),
	}
	goodExpectedResult := &GetArticleResult{
		MTime:   1607209066,
		Content: testContent1,
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
			name: "get_article_success",
			args: args{
				params: goodParams,
				path:   goodPath,
			},
			wantResult: goodExpectedResult,
			wantErr:    false,
		},
		{
			name: "invalid_params",
			args: args{
				params: "invalid_params",
				path:   goodPath,
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "invalid_path",
			args: args{
				params: goodParams,
				path:   "invalid_path",
			},
			wantResult: nil,
			wantErr:    true,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := GetArticleAllGuest(tt.args.remoteAddr, tt.args.params, tt.args.path, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticleAllGuest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("GetArticleAllGuest() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
	wg.Wait()
}
