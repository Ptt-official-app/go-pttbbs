package api

import (
	"reflect"
	"sync"
	"testing"
)

func TestCheckExistsUser(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	params0 := &CheckExistsUserParams{Username: "SYSOP"}
	result0 := &CheckExistsUserResult{UserID: "SYSOP", IsExists: true}

	params1 := &CheckExistsUserParams{Username: "nonexists"}
	result1 := &CheckExistsUserResult{UserID: "", IsExists: false}

	params2 := &CheckExistsUserParams{Username: "non-exists"}

	type args struct {
		remoteAddr string
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
			args:           args{params: params0},
			expectedResult: result0,
		},
		{
			args:           args{params: params1},
			expectedResult: result1,
		},
		{
			args:    args{params: params2},
			wantErr: true,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := CheckExistsUser(tt.args.remoteAddr, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExistsUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("CheckExistsUser() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
