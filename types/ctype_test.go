package types

import (
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Isalpha(tt.args.c); got != tt.expected {
				t.Errorf("Isalpha() = %v, want %v", got, tt.expected)
			}
		})
	}
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Isnumber(tt.args.c); got != tt.expected {
				t.Errorf("Isnumber() = %v, want %v", got, tt.expected)
			}
		})
	}
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Isalnum(tt.args.c); got != tt.expected {
				t.Errorf("Isalnum() = %v, want %v", got, tt.expected)
			}
		})
	}
}
