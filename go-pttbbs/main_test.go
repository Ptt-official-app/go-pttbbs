package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/api"
)

func Test_Index(t *testing.T) {
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
			args: args{path: "/", params: &api.IndexParams{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := initGin()

			w := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.args.params)
			req, _ := http.NewRequest("POST", tt.args.path, bytes.NewBuffer(jsonStr))
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("code: %v", w.Code)
			}
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
