package initgin

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func Test_CreateComment(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testNewPostUserRaw1)
	_ = ptt.SetupNewUser(testUserecRaw3)

	params0 := &api.CreateCommentParams{
		CommentType: ptttype.COMMENT_TYPE_RECOMMEND,
		Content:     []byte("test123"),
	}

	filename0 := &ptttype.Filename_t{}
	copy(filename0[:], []byte("M.1607202239.A.30D"))
	articleID0 := bbs.ToArticleID(filename0)

	type args struct {
		path     string
		username string
		passwd   string
		params   interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			args: args{
				path:     "/board/10_WhoAmI/article/" + string(articleID0) + "/comment",
				username: "A1",
				passwd:   "123123",
				params:   params0,
			},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			router, _ := InitGin()

			jwt := getJwt(router, tt.args.username, tt.args.passwd)
			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, tt.args.params, jwt, nil, "POST")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
		})
		wg.Wait()
	}

}
