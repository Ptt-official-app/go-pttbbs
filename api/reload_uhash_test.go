package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestReloadUHash(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult interface{}
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP"},
			expectedResult: &ReloadUHashResult{Success: true},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := ReloadUHash(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReloadUHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("ReloadUHash() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
