package bbs

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func TestGetArticleAllGuest(t *testing.T) {
	setupTest()
	defer teardownTest()

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))
	articleID := ToArticleID(filename1)

	filename := "testcase/boards/W/WhoAmI/M.1607202239.A.30D"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	logrus.Infof("TestGetArticleAllGuest: start")

	type args struct {
		bboardID   BBoardID
		articleID  ArticleID
		retrieveTS types.Time4
		isHash     bool
	}
	tests := []struct {
		name        string
		args        args
		wantContent []byte
		wantMtime   types.Time4
		wantSum     cmsys.Fnv64_t
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				bboardID:  BBoardID("WhoAmI"),
				articleID: articleID,
			},
			wantContent: testContent1,
			wantMtime:   1607209066,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotContent, gotMtime, gotSum, err := GetArticleAllGuest(tt.args.bboardID, tt.args.articleID, tt.args.retrieveTS, tt.args.isHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticleAllGuest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("GetArticleAllGuest() gotContent = %v, want %v", gotContent, tt.wantContent)
			}
			if !reflect.DeepEqual(gotMtime, tt.wantMtime) {
				t.Errorf("GetArticleAllGuest() gotMtime = %v, want %v", gotMtime, tt.wantMtime)
			}
			if !reflect.DeepEqual(gotSum, tt.wantSum) {
				t.Errorf("GetArticleAllGuest() gotSum = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
	wg.Wait()
}
