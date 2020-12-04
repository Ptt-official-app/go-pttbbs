package cmsys

import (
	"testing"
)

func TestIsLeapYear(t *testing.T) {
	type args struct {
		year int
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{1},
			expected: false,
		},
		{
			args:     args{4},
			expected: true,
		},
		{
			args:     args{25},
			expected: false,
		},
		{
			args:     args{100},
			expected: false,
		},
		{
			args:     args{203},
			expected: false,
		},
		{
			args:     args{240},
			expected: true,
		},
		{
			args:     args{400},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLeapYear(tt.args.year); got != tt.expected {
				t.Errorf("IsLeapYear() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetHoroscope(t *testing.T) {
	type args struct {
		m int
		d int
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		// TODO: Add test cases.
		{
			name:     "01/03: 摩羯",
			args:     args{1, 3},
			expected: 1,
		},
		{
			name:     "01/20: 水瓶",
			args:     args{1, 20},
			expected: 2,
		},
		{
			name:     "02/14: 水瓶",
			args:     args{2, 14},
			expected: 2,
		},
		{
			name:     "12/21: 射手",
			args:     args{12, 21},
			expected: 12,
		},
		{
			name:     "12/22: 摩羯",
			args:     args{12, 22},
			expected: 1,
		},
		{
			name:     "-1/22: 摩羯",
			args:     args{-1, 22},
			expected: 1,
		},
		{
			name:     "14/22: 摩羯",
			args:     args{14, 22},
			expected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHoroscope(tt.args.m, tt.args.d); got != tt.expected {
				t.Errorf("GetHoroscope() = %v, want %v", got, tt.expected)
			}
		})
	}
}
