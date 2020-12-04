package cmsys

import (
	"testing"
)

func Test_fnv32Bytes(t *testing.T) {
	type args struct {
		theBytes []byte
		hval     uint32
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312"), 12345},
			expected: 0x580bef3e,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fnv32Bytes(tt.args.theBytes, tt.args.hval); got != tt.expected {
				t.Errorf("fnv32Bytes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_fnv1a32Bytes(t *testing.T) {
	type args struct {
		theBytes []byte
		hval     uint32
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312"), 12345},
			expected: 0x2d7500e0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fnv1a32Bytes(tt.args.theBytes, tt.args.hval); got != tt.expected {
				t.Errorf("fnv1a32Bytes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_fnv1a32StrCase(t *testing.T) {
	type args struct {
		theBytes []byte
		hval     uint32
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312"), 12345},
			expected: 0x2d7500e0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fnv1a32StrCase(tt.args.theBytes, tt.args.hval); got != tt.expected {
				t.Errorf("fnv1a32StrCase() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_fnv1a32DBCSCase(t *testing.T) {
	type args struct {
		theBytes []byte
		hval     uint32
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312"), 12345},
			expected: 0x2d7500e0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fnv1a32DBCSCase(tt.args.theBytes, tt.args.hval); got != tt.expected {
				t.Errorf("fnv1a32DBCSCase() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_fnv64Bytes(t *testing.T) {
	type args struct {
		theBytes []byte
		hval     uint64
	}
	tests := []struct {
		name     string
		args     args
		expected uint64
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312"), 12345},
			expected: 0x772a7b721d0a27fe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fnv64Bytes(tt.args.theBytes, tt.args.hval); got != tt.expected {
				t.Errorf("fnv64Bytes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_fnv1a64Bytes(t *testing.T) {
	type args struct {
		theBytes []byte
		hval     uint64
	}
	tests := []struct {
		name     string
		args     args
		expected uint64
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312"), 12345},
			expected: 0xf0f6a669273f4180,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fnv1a64Bytes(tt.args.theBytes, tt.args.hval); got != tt.expected {
				t.Errorf("fnv1a64Bytes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_fnv1a64StrCase(t *testing.T) {
	type args struct {
		theBytes []byte
		hval     uint64
	}
	tests := []struct {
		name     string
		args     args
		expected uint64
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312"), 12345},
			expected: 0xf0f6a669273f4180,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fnv1a64StrCase(tt.args.theBytes, tt.args.hval); got != tt.expected {
				t.Errorf("fnv1a64StrCase() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_fnv1a64DBCSCase(t *testing.T) {
	type args struct {
		theBytes []byte
		hval     uint64
	}
	tests := []struct {
		name     string
		args     args
		expected uint64
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312"), 12345},
			expected: 0xf0f6a669273f4180,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fnv1a64DBCSCase(tt.args.theBytes, tt.args.hval); got != tt.expected {
				t.Errorf("fnv1a64DBCSCase() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_fnv1aByte(t *testing.T) {
	type args struct {
		theByte byte
		hval    uint32
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{'a', 12345},
			expected: 0x584c1a88,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fnv1aByte(tt.args.theByte, tt.args.hval); got != tt.expected {
				t.Errorf("fnv1aByte() = %v, want %v", got, tt.expected)
			}
		})
	}
}
