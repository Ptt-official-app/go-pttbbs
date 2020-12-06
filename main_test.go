package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
	log "github.com/sirupsen/logrus"
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
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
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
			args: args{path: "/login", params: &api.LoginParams{UserID: "SYSOP", Passwd: "123123", IP: "127.0.0.1"}}, // json: {}
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

func Test_initAllConfig(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		filename string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantHTTP_POST string
	}{
		// TODO: Add test cases.
		{
			args:          args{"testcase/test.ini"},
			wantHTTP_POST: "test.dev",
		},
		{
			args:    args{"testcase/non-exists.ini"},
			wantErr: true,
		},
		{
			args:    args{"testcase/non-exists"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := initAllConfig(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("initConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(HTTP_HOST, tt.wantHTTP_POST) {
				t.Errorf("initConfig() HTTP_HOST: %v want :%v", HTTP_HOST, tt.wantHTTP_POST)
			}
		})
	}
}

func Test_initMain(t *testing.T) {
	setupTest()
	defer teardownTest()

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initMain(); (err != nil) != tt.wantErr {
				t.Errorf("initMain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
