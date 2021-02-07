package initgin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func Test_CreateBoard(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testUserecRaw3)

	params := &api.CreateBoardParams{
		Brdname: "mnewtest",
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
		// TODO: Add test cases.
		{
			args: args{
				path:     "/class/2/board",
				username: "test",
				passwd:   "123123",
				params:   params,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := InitGin()

			jwt := getJwt(router, tt.args.username, tt.args.passwd)
			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, tt.args.params, jwt, nil, "POST")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
		})
	}

}
