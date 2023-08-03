package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestIsBoardValidUser(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	path0 := &IsBoardValidUserPath{
		BoardID: "10_WhoAmI",
	}

	result0 := &IsBoardValidUserResult{
		IsValid: true,
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
		path       interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult interface{}
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", path: path0},
			expectedResult: result0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := IsBoardValidUser(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsBoardValidUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("IsBoardValidUser() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
