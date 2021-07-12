package main

import (
	"reflect"
	"sync"
	"testing"
)

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
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
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
	wg.Wait()
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
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := initMain(); (err != nil) != tt.wantErr {
				t.Errorf("initMain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	wg.Wait()
}
