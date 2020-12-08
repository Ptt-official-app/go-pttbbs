package api

import (
	"reflect"
	"testing"
)

func TestIndex(t *testing.T) {
	type args struct {
		userID string
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
			args:     args{"SYSOP", nil},
			expected: &IndexResult{Data: "index"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			got, err := Index(tt.args.userID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Index() = %v, want %v", got, tt.expected)
			}
		})
	}
}
