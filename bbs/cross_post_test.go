package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestCrossPost(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testUserecRaw3)

	boardSummary, _ := CreateBoard("test", 5, "newboard0", []byte("CPBL"), []byte("new-board"), nil, 0, 0, 0, false)

	forwardBoardSummary, _ := CreateBoard("test", 5, "fwboard0", []byte("CPBL"), []byte("fw-board"), nil, 0, 0, 0, false)

	class0 := []byte("test")
	title0 := []byte("this is a test")
	content0 := [][]byte{[]byte("test1"), []byte("test2")}
	ip0 := "127.0.0.1"

	fwTitle0 := []byte("Fw: [test] this is a test")

	articleSummary0, err := CreateArticle("test", boardSummary.BBoardID, class0, title0, content0, ip0)
	if err != nil {
		t.Errorf("unable to create article: e: %v", err)
	}

	expectedClass0 := ptttype.ARTICLE_CLASS_FORWARD
	expectedSummary0 := &ArticleSummary{
		BBoardID: "14_fwboard0",
		Owner:    "SYSOP",
		Title:    fwTitle0,
		Class:    expectedClass0,
		Idx:      "",
	}

	expectedComment0 := []byte{
		0xa1, 0xb0, 0x20, // ※
		0x1b, 0x5b, 0x31, 0x3b, 0x33, 0x32, 0x6d, // [1;32m
		0x53, 0x59, 0x53, 0x4f, 0x50, // SYSOP
		0x1b, 0x5b, 0x30, 0x3b, 0x33, 0x32, 0x6d, // [0;32m:
		0x3a, 0xc2, 0xe0, 0xbf, 0xfd, 0xa6, 0xdc, 0xac, 0xdd, 0xaa, 0x4f, 0x20, // 轉錄至看板
		0x66, 0x77, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x30, // fwboard0
		0x1b, 0x5b, 0x6d, //
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, //
		0x00, 0x00, 0x2f, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x0a, // 00/00 00:00

	}

	removeCommentIdxes := []int{
		84, 85, 87, 88, 90, 91, 93, 94,
	}

	type args struct {
		uuserID   UUserID
		bboardID  BBoardID
		articleID ArticleID
		xboardID  BBoardID
		ip        string
	}
	tests := []struct {
		name                   string
		args                   args
		expectedArticleSummary *ArticleSummary
		expectedComment        []byte
		expectedCommentMTime   types.Time4
		wantErr                bool
	}{
		// TODO: Add test cases.
		{
			args:                   args{uuserID: "SYSOP", bboardID: boardSummary.BBoardID, articleID: articleSummary0.ArticleID, xboardID: forwardBoardSummary.BBoardID, ip: "127.0.0.1"},
			expectedArticleSummary: expectedSummary0,
			expectedComment:        expectedComment0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotArticleSummary, gotComment, gotCommentMTime, err := CrossPost(tt.args.uuserID, tt.args.bboardID, tt.args.articleID, tt.args.xboardID, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrossPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotArticleSummary.Filename = ""
			gotArticleSummary.CreateTime = 0
			gotArticleSummary.ArticleID = ""
			gotArticleSummary.Idx = ""
			testutil.TDeepEqual(t, "article", gotArticleSummary, tt.expectedArticleSummary)

			for _, idx := range removeCommentIdxes {
				if idx >= len(gotComment) {
					break
				}
				gotComment[idx] = 0x00
			}

			testutil.TDeepEqual(t, "comment", gotComment, tt.expectedComment)

			if gotCommentMTime == 0 {
				t.Errorf("CrossPost() gotCommentMTime == 0")
			}
		})
		wg.Wait()
	}
}
