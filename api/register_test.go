package api

import (
	"sync"
	"testing"
)

func TestRegister(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		params interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args: args{params: &RegisterParams{Username: "C1", Passwd: "567"}},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := Register(testIP, tt.args.params, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result, ok := got.(*RegisterResult)
			if !ok {
				t.Errorf("Register(), not register-result")
				return
			}
			jwt := result.Jwt
			if len(jwt) == 0 {
				t.Errorf("Register() no jwt")
				return
			}
		})
		wg.Wait()
	}
}
