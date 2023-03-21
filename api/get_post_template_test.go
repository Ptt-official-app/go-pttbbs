package api

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestGetPostTemplate(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	// keep chtime for testing
	filename := "testcase/boards/W/WhoAmI/postsample.0"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	goodUUserID := bbs.UUserID("SYSOP")
	goodParams := &GetPostTemplateParams{}
	goodPath := &GetPostTemplatePath{
		BBoardID:   bbs.BBoardID("10_WhoAmI"),
		TemplateID: 1,
	}
	goodExpectedResult := &GetPostTemplateResult{
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
		name       string
		args       args
		wantResult interface{}
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "get_post_template",
			args: args{
				uuserID: goodUUserID,
				params:  goodParams,
				path:    goodPath,
			},
			wantResult: goodExpectedResult,
			wantErr:    false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := GetPostTemplate(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("GetPostTemplate() = %v, want %v", gotResult, tt.wantResult)
			}
		})
		wg.Wait()
	}
}
