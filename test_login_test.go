package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
	log "github.com/sirupsen/logrus"
)

func Test_Login(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		path   string
		params interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			args: args{path: withPrefix("/token"), params: &api.LoginParams{UserID: "SYSOP", Passwd: "123123", IP: "127.0.0.1"}}, // json: {}
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Infof("Test_Login: start: name: %v", tt.name)

			router, _ := initGin()

			w := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.args.params)
			req, _ := http.NewRequest("POST", tt.args.path, bytes.NewBuffer(jsonStr))
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}

			body, _ := ioutil.ReadAll(w.Body)

			results := &api.LoginResult{}
			json.Unmarshal(body, results)

			jwt := results.Jwt

			log.Infof("jwt: %v", jwt)

		})
	}
}
