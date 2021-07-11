package api

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestCreateComment(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testNewPostUserRaw1)

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))
	articleID := bbs.ToArticleID(filename1)

	params0 := &CreateCommentParams{
		CommentType: ptttype.COMMENT_TYPE_RECOMMEND,
		Content:     []byte("test123"),
	}
	path0 := &CreateCommentPath{
		BBoardID:  "10_WhoAmI",
		ArticleID: articleID,
	}
	expected0 := &CreateCommentResult{
		Content: []byte("\x1b[1;37m\xb1\xc0 \x1b[33mA1\x1b[m\x1b[33m: test123                                                     \x1b[m 05/26 10:25\n"),
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
		expectedResult *CreateCommentResult
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "A1", params: params0, path: path0},
			expectedResult: expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := CreateComment(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := gotResult.(*CreateCommentResult)
			got.MTime = 0
			theLen := len(got.Content) - 12
			testutil.TDeepEqual(t, "got", got.Content[:theLen], tt.expectedResult.Content[:theLen])
		})
	}
}
