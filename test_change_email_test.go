package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
)

func Test_ChangeEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	jwt, _ := api.CreateEmailToken("SYSOP", "", "test@ptt.test", api.CONTEXT_CHANGE_EMAIL)

	params0 := &api.ChangeEmailParams{
		Jwt: jwt,
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
				path:     "/user/SYSOP/changeemail",
				username: "SYSOP",
				passwd:   "123123",
				params:   params0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := initGin()

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
