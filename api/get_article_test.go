package api

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestGetArticle(t *testing.T) {
	setupTest()
	defer teardownTest()

	params1 := &GetArticleParams{}
	path1 := &GetArticlePath{
		BBoardID:  bbs.BBoardID("10_WhoAmI"),
		ArticleID: bbs.ArticleID("1Vo_M_CDSYSOP"),
	}

	filename := "testcase/boards/W/WhoAmI/M.1607202239.A.30D"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	expectedResult1 := &GetArticleResult{
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
		// TODO: Add test cases.
		{
			args: args{
				uuserID: bbs.UUserID("SYSOP"),
				params:  params1,
				path:    path1,
			},
			expectedResult: expectedResult1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
}
