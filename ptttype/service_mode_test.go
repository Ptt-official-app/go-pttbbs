package ptttype

import (
	"sync"
	"testing"
)

func TestServiceMode_String(t *testing.T) {
	tests := []struct {
		name     string
		s        ServiceMode
		expected string
	}{
		// TODO: Add test cases.
		{
			s:        DEV,
			expected: "DEV",
		},
		{
			s:        PRODUCTION,
			expected: "PRODUCTION",
		},
		{
			s:        DEBUG,
			expected: "DEBUG",
		},
		{
			s:        ServiceMode(100),
			expected: "UNKNOWN",
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := tt.s.String(); got != tt.expected {
				t.Errorf("ServiceMode.String() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func Test_stringToServiceMode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name     string
		args     args
		expected ServiceMode
	}{
		// TODO: Add test cases.
		{
			args:     args{"DEV"},
			expected: DEV,
		},
		{
			args:     args{"PRODUCTION"},
			expected: PRODUCTION,
		},
		{
			args:     args{"DEBUG"},
			expected: DEBUG,
		},
		{
			args:     args{"UNKNOWN"},
			expected: DEV,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := stringToServiceMode(tt.args.str); got != tt.expected {
				t.Errorf("stringToServiceMode() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}
