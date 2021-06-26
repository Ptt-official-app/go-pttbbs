package bbs

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/sirupsen/logrus"
)

func TestCreateComment(t *testing.T) {
	setupTest()
	defer teardownTest()

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))
	articleID := ToArticleID(filename1)

	content0 := []byte("test123")
	expected0 := []byte("\x1b[1;37m\xb1\xc0 \x1b[33mCodingMan\x1b[m\x1b[33m: test123                                              \x1b[m 05/26 10:25\n")

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotComment, _, err := CreateComment(tt.args.userID, tt.args.boardID, tt.args.articleID, tt.args.commentType, tt.args.content, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRecommend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logrus.Infof("gotComment: %v", gotComment)
			theLen := len(gotComment) - 12
			if !reflect.DeepEqual(gotComment[:theLen], tt.expectedComment[:theLen]) {
				t.Errorf("CreateRecommend() = %v, want %v", gotComment, tt.expectedComment)
			}
		})
	}
}
