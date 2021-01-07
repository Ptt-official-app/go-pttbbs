package main

import (
	"net/http"
	"net/http/httptest"
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
				path:     "/users/SYSOP/attemptchangeemail",
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
