package types

import (
	"sync"
	"testing"
)

func TestIsalpha(t *testing.T) {
	type args struct {
		c byte
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{' '},
			expected: false,
		},
		{
			args:     args{'/'},
			expected: false,
		},
		{
			args:     args{'0'},
			expected: false,
		},
		{
			args:     args{'9'},
			expected: false,
		},
		{
			args:     args{'A'},
			expected: true,
		},
		{
			args:     args{'Z'},
			expected: true,
		},
		{
			args:     args{'a'},
			expected: true,
		},
		{
			args:     args{'z'},
			expected: true,
		},
		{
			args:     args{'~'},
			expected: false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := Isalpha(tt.args.c); got != tt.expected {
				t.Errorf("Isalpha() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestIsnumber(t *testing.T) {
	type args struct {
		c byte
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{' '},
			expected: false,
		},
		{
			args:     args{'/'},
			expected: false,
		},
		{
			args:     args{'0'},
			expected: true,
		},
		{
			args:     args{'9'},
			expected: true,
		},
		{
			args:     args{'A'},
			expected: false,
		},
		{
			args:     args{'Z'},
			expected: false,
		},
		{
			args:     args{'a'},
			expected: false,
		},
		{
			args:     args{'z'},
			expected: false,
		},
		{
			args:     args{'~'},
			expected: false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := Isnumber(tt.args.c); got != tt.expected {
				t.Errorf("Isnumber() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestIsalnum(t *testing.T) {
	type args struct {
		c byte
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{' '},
			expected: false,
		},
		{
			args:     args{'/'},
			expected: false,
		},
		{
			args:     args{'0'},
			expected: true,
		},
		{
			args:     args{'9'},
			expected: true,
		},
		{
			args:     args{'A'},
			expected: true,
		},
		{
			args:     args{'Z'},
			expected: true,
		},
		{
			args:     args{'a'},
			expected: true,
		},
		{
			args:     args{'z'},
			expected: true,
		},
		{
			args:     args{'~'},
			expected: false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := Isalnum(tt.args.c); got != tt.expected {
				t.Errorf("Isalnum() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}
