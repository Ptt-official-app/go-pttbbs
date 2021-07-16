package initgin

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
)

func Test_AttemptChangeEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	//until Fri Nov 14 01:28:37 EST 2245
	params0 := &api.AttemptChangeEmailParams{
		ClientInfo: "test_clientinfo",
		Passwd:     "123123",
		Email:      "test@ptt.test",
	}

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
				path:     "/user/SYSOP/attemptchangeemail",
				username: "SYSOP",
				passwd:   "123123",
				params:   params0,
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
			req := setRequest(tt.args.path, tt.args.params, jwt, nil, "POST")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
		})
		wg.Wait()
	}
}
