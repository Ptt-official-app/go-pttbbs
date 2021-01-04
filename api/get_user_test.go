package api

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestGetUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	path := &GetUserPath{UserID: "SYSOP"}
	expected := GetUserResult(testUserec)

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
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", path: path},
			expectedResult: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := GetUser(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := gotResult.(GetUserResult)
			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("GetUser() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
	}
}
