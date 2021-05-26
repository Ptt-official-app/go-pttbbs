package path

import "testing"

func TestSetDIRPath(t *testing.T) {
	type args struct {
		dirFilename string
		filename    string
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		// TODO: Add test cases.
		{
			args:     args{dirFilename: "boards/W/WhoAmI/.DIR", filename: "test"},
			expected: "boards/W/WhoAmI/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetDIRPath(tt.args.dirFilename, tt.args.filename); got != tt.expected {
				t.Errorf("SetDIRPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}
