package cmsys

import (
	"testing"
)

func TestStringHashWithHashBits(t *testing.T) {
	type args struct {
		theBytes []byte
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312")},
			expected: 0x94b2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringHashWithHashBits(tt.args.theBytes); got != tt.expected {
				t.Errorf("StringHashWithHashBits() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringHash(t *testing.T) {
	type args struct {
		theBytes []byte
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312")},
			expected: 0x7aeb94b2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringHash(tt.args.theBytes); got != tt.expected {
				t.Errorf("StringHash() = %v, want %v", got, tt.expected)
			}
		})
	}
}
