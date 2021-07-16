package ptt

import (
	"os"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/sirupsen/logrus"
)

func TestRecommend(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	boardID := &ptttype.BoardID_t{}
	copy(boardID[:], "WhoAmI")
	filename0 := &ptttype.Filename_t{}
	copy(filename0[:], "M.1607203395.A.00D")
	expected0 := []byte("\x1b[1;37m\xb1\xc0 \x1b[33mA1\x1b[m\x1b[33m: test123                                                     \x1b[m 05/26 10:25\n")
	expected1 := []byte("\x1b[1;31m\xa1\xf7 \x1b[33mA1\x1b[m\x1b[33m:test456                                                   \x1b[m\xb1\xc0 05/27 02:31\n")

	type args struct {
		user        *ptttype.UserecRaw
		uid         ptttype.UID
		boardID     *ptttype.BoardID_t
		bid         ptttype.Bid
		filename    *ptttype.Filename_t
		commentType ptttype.CommentType
		content     []byte
		ip          *ptttype.IPv4_t
		from        []byte
	}
	tests := []struct {
		name            string
		args            args
		isOldRecommend  bool
		isSmartMerge    bool
		expectedComment []byte
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{user: testNewPostUser1, uid: 2, boardID: boardID, bid: 10, filename: filename0, commentType: ptttype.COMMENT_TYPE_RECOMMEND, content: []byte("test123")},
			expectedComment: expected0,
		},
		{
			args:            args{user: testNewPostUser1, uid: 2, boardID: boardID, bid: 10, filename: filename0, commentType: ptttype.COMMENT_TYPE_RECOMMEND, content: []byte("test456")},
			isOldRecommend:  true,
			isSmartMerge:    true,
			expectedComment: expected1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			total, _ := cache.GetBTotalWithRetry(tt.args.bid)

			dirFilename, _ := setBDir(boardID)
			_, origFhdr, _ := cmsys.GetRecord(dirFilename, tt.args.filename, int(total))

			origOldRecommend := ptttype.OLDRECOMMEND
			origSmartMerge := ptttype.EDITPOST_SMARTMERGE
			defer func() {
				ptttype.OLDRECOMMEND = origOldRecommend
				ptttype.EDITPOST_SMARTMERGE = origSmartMerge
			}()

			ptttype.OLDRECOMMEND = tt.isOldRecommend
			ptttype.EDITPOST_SMARTMERGE = tt.isSmartMerge

			gotComment, _, err := Recommend(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, tt.args.filename, tt.args.commentType, tt.args.content, tt.args.ip, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("Recommend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			logrus.Infof("gotComment: %v", string(gotComment))

			theLen := len(tt.expectedComment) - 12

			testutil.TDeepEqual(t, "got", gotComment[:theLen], tt.expectedComment[:theLen])

			filename, _ := path.SetBFile(tt.args.boardID, tt.args.filename.String())
			file, _ := os.Open(filename)
			defer file.Close()

			theBytes := make([]byte, 500)
			n, _ := file.Read(theBytes)
			theBytes = theBytes[:n]
			startIdx := len(theBytes) - len(tt.expectedComment)
			testutil.TDeepEqual(t, "comment", theBytes[startIdx:], gotComment)

			_, fhdr, _ := cmsys.GetRecord(dirFilename, tt.args.filename, int(total))
			if fhdr.Recommend <= origFhdr.Recommend {
				t.Errorf("Recommend: recommend: %v orig: %v", fhdr.Recommend, origFhdr.Recommend)
			}
		})
		wg.Wait()
	}
}
