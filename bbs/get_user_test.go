package bbs

import (
	"reflect"
	"testing"
)

func TestGetUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID UUserID
	}
	tests := []struct {
		name         string
		args         args
		expectedUser *Userec
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			args:         args{uuserID: "SYSOP"},
			expectedUser: testUserec1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := GetUser(tt.args.uuserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.expectedUser) {
				t.Errorf("GetUser() = %v, want %v", gotUser, tt.expectedUser)
			}
		})
	}
}
