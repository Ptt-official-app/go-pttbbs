package initgin

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
)

func Test_GetArticle(t *testing.T) {
	setupTest()
	defer teardownTest()

	params1 := &api.GetArticleParams{}

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
		{
			args: args{
				path:     "/board/10_WhoAmI/article/1Vo_M_CDSYSOP",
				username: "SYSOP",
				passwd:   "123123",
				params:   params1,
			},
		},
	}

	router, _ := InitGin()
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			jwt := getJwt(router, tt.args.username, tt.args.passwd)
			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, tt.args.params, jwt, nil, "GET")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
		})
		wg.Wait()
	}
}
