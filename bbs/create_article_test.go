package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestCreateArticle(t *testing.T) {
	// setupTest in for-loop
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testNewPostUserRaw1)

	class0 := []byte("test")
	title0 := []byte("this is a test")
	fullTitle0 := []byte("[test] this is a test")
	content0 := [][]byte{[]byte("test1"), []byte("test2")}
	ip0 := "127.0.0.1"

	expectedSummary0 := &ArticleSummary{
		BBoardID:  "10_WhoAmI",
		Owner:     "A1",
		FullTitle: fullTitle0,
		Class:     class0,
		Idx:       "",
		RealTitle: []byte("this is a test"),
	}

	type args struct {
		uuserID  UUserID
		bboardID BBoardID
		posttype []byte
		title    []byte
		content  [][]byte
		ip       string
	}
	tests := []struct {
		name            string
		args            args
		expectedSummary *ArticleSummary
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				uuserID:  "A1",
				bboardID: "10_WhoAmI",
				posttype: class0,
				title:    title0,
				content:  content0,
				ip:       ip0,
			},
			expectedSummary: expectedSummary0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, err := CreateArticle(tt.args.uuserID, tt.args.bboardID, tt.args.posttype, tt.args.title, tt.args.content, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotSummary.ArticleID = ""
			gotSummary.Filename = ""
			gotSummary.CreateTime = 0
			gotSummary.MTime = 0
			gotSummary.Idx = ""
			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)
		})
		wg.Wait()
	}
}
