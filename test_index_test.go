package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
)

func Test_Index(t *testing.T) {
	setupTest()
	defer teardownTest()

	params := &LoginRequiredParams{
		Data: &api.IndexParams{},
	}
	type args struct {
		path   string
		userID string
		passwd string
		params interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			args: args{
				path:   "/",
				userID: "SYSOP",
				passwd: "123123",
				params: params,
			}, // json: {"Data": "index"}
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := initGin()

			w := httptest.NewRecorder()
			loginParams := &api.LoginParams{UserID: tt.args.userID, Passwd: tt.args.passwd, IP: "127.0.0.1"}
			jsonStr, _ := json.Marshal(loginParams)
			req, _ := http.NewRequest("POST", withPrefix("/token"), bytes.NewBuffer(jsonStr))
			router.ServeHTTP(w, req)
			body, _ := ioutil.ReadAll(w.Body)
			resultLogin := &api.LoginResult{}
			_ = json.Unmarshal(body, resultLogin)
			jwt := resultLogin.Jwt

			params.UserID = tt.args.userID
			params.Jwt = jwt

			w = httptest.NewRecorder()
			jsonStr, _ = json.Marshal(tt.args.params)
			req, _ = http.NewRequest("POST", tt.args.path, bytes.NewBuffer(jsonStr))
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
		})
	}
}
