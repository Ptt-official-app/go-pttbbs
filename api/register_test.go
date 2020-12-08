package api

import (
	"testing"
)

func TestRegister(t *testing.T) {
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
			args: args{params: &RegisterParams{UserID: "C1", Passwd: "567"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			got, err := Register(testIP, tt.args.params)
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
	}
}
