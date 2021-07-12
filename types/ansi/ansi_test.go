package ansi

import (
	"sync"
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
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := ANSIColor(tt.args.color); got != tt.expected {
				t.Errorf("ANSIColor() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
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
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := ANSIReset(); got != tt.expected {
				t.Errorf("ANSIReset() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}
