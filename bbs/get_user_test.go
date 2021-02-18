package bbs

import (
	"reflect"
	"sync"
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
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
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
	wg.Wait()
}
