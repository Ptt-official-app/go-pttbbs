package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_IsBoardValidUser(t *testing.T) {
	setupTest()
	defer teardownTest()

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
				path:     "/board/10_WhoAmI/isvalid",
				username: "SYSOP",
				passwd:   "123123",
				params:   nil,
			}, // json: {}
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := initGin()

			jwt := getJwt(router, tt.args.username, tt.args.passwd)
			w := httptest.NewRecorder()
			req := setRequest(tt.args.path, tt.args.params, jwt, nil, "GET")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}

		})
	}

}
