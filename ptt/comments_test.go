package ptt

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/sirupsen/logrus"
)

func TestFormatCommentString(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	testBoard10, _ := cache.GetBCache(10)
	testIP := &ptttype.IPv4_t{}
	copy(testIP[:], []byte("192.168.0.1"))
	expected0 := []byte("\x1b[1;37m\xb1\xc0 \x1b[33mCodingMan\x1b[m\x1b[33m: test123                                              \x1b[m 05/26 10:25\n")
	expected1 := []byte("\x1b[1;31m\xbcN \x1b[33mCodingMan\x1b[m\x1b[33m: test123                                              \x1b[m 05/26 10:25\n")
	expected2 := []byte("\x1b[1;31m\xa1\xf7 \x1b[33mCodingMan\x1b[m\x1b[33m: test123                                              \x1b[m 05/26 10:25\n")

	type args struct {
		user        *ptttype.UserecRaw
		board       *ptttype.BoardHeaderRaw
		commentType ptttype.CommentType
		content     []byte
		ip          *ptttype.IPv4_t
		from        []byte
	}
	tests := []struct {
		name            string
		args            args
		expectedComment []byte
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:        testUserecRaw2,
				board:       testBoard10,
				commentType: ptttype.COMMENT_TYPE_RECOMMEND,
				content:     []byte("test123"),
				ip:          testIP,
			},
			expectedComment: expected0,
		},
		{
			args: args{
				user:        testUserecRaw2,
				board:       testBoard10,
				commentType: ptttype.COMMENT_TYPE_BOO,
				content:     []byte("test123"),
				ip:          testIP,
			},
			expectedComment: expected1,
		},
		{
			args: args{
				user:        testUserecRaw2,
				board:       testBoard10,
				commentType: ptttype.COMMENT_TYPE_COMMENT,
				content:     []byte("test123"),
				ip:          testIP,
			},
			expectedComment: expected2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotComment, err := FormatCommentString(tt.args.user, tt.args.board, tt.args.commentType, tt.args.content, tt.args.ip, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatCommentString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logrus.Infof("gotComment: %v", string(gotComment))
			theLen := len(tt.expectedComment) - 12
			testutil.TDeepEqual(t, "got", gotComment[:theLen], tt.expectedComment[:theLen])
		})
	}
}
