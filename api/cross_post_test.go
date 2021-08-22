package api

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestCrossPost(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	_ = ptt.SetupNewUser(testUserecRaw3)

	boardSummary, _ := bbs.CreateBoard("test", 5, "newboard0", []byte("CPBL"), []byte("new-board"), nil, 0, 0, 0, false)

	forwardBoardSummary, _ := bbs.CreateBoard("test", 5, "fwboard0", []byte("CPBL"), []byte("fw-board"), nil, 0, 0, 0, false)

	class0 := []byte("test")
	title0 := []byte("this is a test")
	content0 := [][]byte{[]byte("test1"), []byte("test2")}
	ip0 := "127.0.0.1"

	fwTitle0 := []byte("Fw: [test] this is a test")

	articleSummary0, err := bbs.CreateArticle("test", boardSummary.BBoardID, class0, title0, content0, ip0)
	if err != nil {
		t.Errorf("unable to create article: e: %v", err)
	}

	expectedSummary0 := &bbs.ArticleSummary{
		BBoardID:    forwardBoardSummary.BBoardID,
		Owner:       "SYSOP",
		FullTitle:   fwTitle0,
		Class:       class0,
		Idx:         "",
		RealTitle:   title0,
		SubjectType: ptttype.SUBJECT_FORWARD,
	}

	params0 := &CrossPostParams{
		XBoardID: forwardBoardSummary.BBoardID,
	}
	path0 := &CrossPostPath{
		BBoardID:  boardSummary.BBoardID,
		ArticleID: articleSummary0.ArticleID,
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

	expected0 := &CrossPostResult{
		ArticleSummary: expectedSummary0,
		Comment:        expectedComment0,
	}

	removeCommentIdxes := []int{
		84, 85, 87, 88, 90, 91, 93, 94,
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
		expectedResult *CrossPostResult
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", params: params0, path: path0},
			expectedResult: expected0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := CrossPost(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrossPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result, _ := gotResult.(*CrossPostResult)

			result.ArticleSummary.Filename = ""
			result.ArticleSummary.CreateTime = 0
			result.ArticleSummary.ArticleID = ""
			result.ArticleSummary.Idx = ""

			result.CommentMTime = 0

			for _, idx := range removeCommentIdxes {
				if idx >= len(result.Comment) {
					break
				}
				result.Comment[idx] = 0x00
			}

			testutil.TDeepEqual(t, "result", result, tt.expectedResult)
		})
		wg.Wait()
	}
}
