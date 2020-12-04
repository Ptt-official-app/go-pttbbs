package ansi

import (
	"testing"
)

func TestANSIColor(t *testing.T) {
	type args struct {
		color string
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		// TODO: Add test cases.
		{
			args:     args{"1;30"},
			expected: "\x1b[1;30m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ANSIColor(tt.args.color); got != tt.expected {
				t.Errorf("ANSIColor() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestANSIReset(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		// TODO: Add test cases.
		{
			expected: "\x1b[m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ANSIReset(); got != tt.expected {
				t.Errorf("ANSIReset() = %v, want %v", got, tt.expected)
			}
		})
	}
}
