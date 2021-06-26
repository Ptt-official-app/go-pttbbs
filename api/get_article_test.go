package api

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestGetArticle(t *testing.T) {
	setupTest()
	defer teardownTest()

	// keep chtime for testing
	filename := "testcase/boards/W/WhoAmI/M.1607202239.A.30D"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	goodUUserID := bbs.UUserID("SYSOP")
	goodParams := &GetArticleParams{}
	goodPath := &GetArticlePath{
		BBoardID:  bbs.BBoardID("10_WhoAmI"),
		ArticleID: bbs.ArticleID("1Vo_M_CD"),
	}
	goodExpectedResult := &GetArticleResult{
		MTime:   1607209066,
		Content: testContent1,
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
		{
			name: "get_article_success",
			args: args{
				uuserID: goodUUserID,
				params:  goodParams,
				path:    goodPath,
			},
			expectedResult: goodExpectedResult,
			wantErr:        false,
		},
		{
			name: "invalid_params",
			args: args{
				uuserID: goodUUserID,
				params:  "invalid_params",
				path:    goodPath,
			},
			expectedResult: nil,
			wantErr:        true,
		},
		{
			name: "invalid_path",
			args: args{
				uuserID: goodUUserID,
				params:  goodParams,
				path:    "invalid_path",
			},
			expectedResult: nil,
			wantErr:        true,
		},
		{
			name: "invalid_get_article",
			args: args{
				uuserID: bbs.UUserID("invalid_user"),
				params:  goodParams,
				path:    goodPath,
			},
			expectedResult: nil,
			wantErr:        true,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := GetArticle(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("GetArticle() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
	}
	wg.Wait()
}
