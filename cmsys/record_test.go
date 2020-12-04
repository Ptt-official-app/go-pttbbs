package cmsys

import "testing"

func TestGetNumRecords(t *testing.T) {
	type args struct {
		filename string
		size     uintptr
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNumRecords(tt.args.filename, tt.args.size); got != tt.expected {
				t.Errorf("GetNumRecords() = %v, want %v", got, tt.expected)
			}
		})
	}
}
