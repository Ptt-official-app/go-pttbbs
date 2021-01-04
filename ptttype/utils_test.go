package ptttype

import "testing"

func TestValidUSHMEntry(t *testing.T) {
	type args struct {
		x UtmpID
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{x: -1},
			expected: false,
		},
		{
			args:     args{x: 0},
			expected: true,
		},
		{
			args:     args{x: USHM_SIZE - 1},
			expected: true,
		},
		{
			args:     args{x: USHM_SIZE},
			expected: false,
		},
		{
			args:     args{x: USHM_SIZE + 1},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidUSHMEntry(tt.args.x); got != tt.expected {
				t.Errorf("ValidUSHMEntry() = %v, want %v", got, tt.expected)
			}
		})
	}
}
