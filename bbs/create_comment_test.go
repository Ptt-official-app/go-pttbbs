package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/sirupsen/logrus"
)

func TestCreateComment(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testNewPostUserRaw1)

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))
	articleID := ToArticleID(filename1)

	content0 := []byte("test123")
	expected0 := []byte("\x1b[1;37m\xb1\xc0 \x1b[33mCodingMan\x1b[m\x1b[33m: test123                                              \x1b[m 05/26 10:25\n")
	expected1 := []byte("\x1b[1;37m\xb1\xc0 \x1b[33mA1\x1b[m\x1b[33m: test123                                                     \x1b[m 05/26 10:25\n")

	type args struct {
		userID      UUserID
		boardID     BBoardID
		articleID   ArticleID
		commentType ptttype.CommentType
		content     []byte
		ip          string
	}
	tests := []struct {
		name            string
		args            args
		expectedComment []byte
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{userID: "CodingMan", boardID: "10_WhoAmI", articleID: articleID, commentType: ptttype.COMMENT_TYPE_RECOMMEND, content: content0, ip: "127.0.0.1"},
			expectedComment: expected0,
		},
		{
			args:            args{userID: "A1", boardID: "10_WhoAmI", articleID: articleID, commentType: ptttype.COMMENT_TYPE_RECOMMEND, content: content0, ip: "127.0.0.1"},
			expectedComment: expected1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotComment, _, err := CreateComment(tt.args.userID, tt.args.boardID, tt.args.articleID, tt.args.commentType, tt.args.content, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRecommend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			logrus.Infof("gotComment: %v, e: %v", gotComment, err)
			theLen := len(gotComment) - 12
			testutil.TDeepEqual(t, "got", gotComment[:theLen], tt.expectedComment[:theLen])
		})
		wg.Wait()
	}
}
