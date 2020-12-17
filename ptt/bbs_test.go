package ptt

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestReadPost(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("WhoAmI"))

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))

	filename := "testcase/boards/W/WhoAmI/M.1607202239.A.30D"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	filename2 := &ptttype.Filename_t{}
	copy(filename2[:], []byte("M.1607202239.A.31D"))
	type args struct {
		user       *ptttype.UserecRaw
		uid        ptttype.Uid
		boardID    *ptttype.BoardID_t
		bid        ptttype.Bid
		filename   *ptttype.Filename_t
		retrieveTS types.Time4
	}
	tests := []struct {
		name            string
		args            args
		expectedContent []byte
		expectedMtime   types.Time4
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:     testUserecRaw1,
				uid:      1,
				boardID:  boardID1,
				bid:      10,
				filename: filename1,
			},
			expectedContent: testContent1,
			expectedMtime:   1607209066,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID1,
				bid:        10,
				filename:   filename1,
				retrieveTS: 1607209066,
			},
			expectedContent: nil,
			expectedMtime:   1607209066,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID1,
				bid:        10,
				filename:   filename2,
				retrieveTS: 1607209066,
			},
			expectedContent: nil,
			expectedMtime:   0,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, gotMtime, err := ReadPost(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, tt.args.filename, tt.args.retrieveTS)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.expectedContent) {
				t.Errorf("ReadPost() gotContent = %v, want %v", gotContent, tt.expectedContent)
			}
			if !reflect.DeepEqual(gotMtime, tt.expectedMtime) {
				t.Errorf("ReadPost() gotMtime = %v, want %v", gotMtime, tt.expectedMtime)
			}
		})
	}
}
