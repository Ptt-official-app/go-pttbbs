package bbs

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestGetArticle(t *testing.T) {
	setupTest()
	defer teardownTest()

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))
	articleID := ToArticleID(filename1)

	filename := "testcase/boards/W/WhoAmI/M.1607202239.A.30D"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	type args struct {
		uuserID    UUserID
		bboardID   BBoardID
		articleID  ArticleID
		retrieveTS types.Time4
		isHash     bool
	}
	tests := []struct {
		name            string
		args            args
		expectedContent []byte
		expectedMtime   types.Time4
		expectedSum     cmsys.Fnv64_t
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				uuserID:   UUserID("SYSOP"),
				bboardID:  BBoardID("10_WhoAmI"),
				articleID: articleID,
			},
			expectedContent: testContent1,
			expectedMtime:   1607209066,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotContent, gotMtime, gotsum, err := GetArticle(tt.args.uuserID, tt.args.bboardID, tt.args.articleID, tt.args.retrieveTS, tt.args.isHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotContent, tt.expectedContent)

			testutil.TDeepEqual(t, "sum", gotsum, tt.expectedSum)

			testutil.TDeepEqual(t, "mtime", gotMtime, tt.expectedMtime)
		})
	}
	wg.Wait()
}
