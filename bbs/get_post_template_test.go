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

func TestGetPostTemplate(t *testing.T) {
	setupTest()
	defer teardownTest()

	filename := "testcase/boards/W/WhoAmI/postsample.0"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	type args struct {
		uuserID    UUserID
		bboardID   BBoardID
		templateID ptttype.SortIdx
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
				uuserID:    "SYSOP",
				bboardID:   "10_WhoAmI",
				templateID: 1,
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
			gotContent, gotMtime, gotSum, err := GetPostTemplate(tt.args.uuserID, tt.args.bboardID, tt.args.templateID, tt.args.retrieveTS, tt.args.isHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotContent, tt.wantContent)
			testutil.TDeepEqual(t, "mtime", gotMtime, tt.wantMtime)
			testutil.TDeepEqual(t, "checksum", gotSum, tt.wantSum)
		})
	}
	wg.Wait()
}
