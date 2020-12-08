package ptttype

import (
	"testing"
)

func TestPERM_HasUserPerm(t *testing.T) {
	type args struct {
		perm PERM
	}
	tests := []struct {
		name     string
		p        PERM
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			p:        PERM_BASIC,
			args:     args{perm: PERM_POST},
			expected: false,
		},
		{
			p:        PERM_BASIC,
			args:     args{perm: PERM_BASIC},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.HasUserPerm(tt.args.perm); got != tt.expected {
				t.Errorf("PERM.HasUserPerm() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPERM_HasBasicUserPerm(t *testing.T) {
	type args struct {
		perm PERM
	}
	tests := []struct {
		name     string
		p        PERM
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			p:        PERM_POST | PERM_BASIC,
			args:     args{perm: PERM_POST},
			expected: true,
		},
		{
			p:        PERM_POST,
			args:     args{perm: PERM_POST},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.HasBasicUserPerm(tt.args.perm); got != tt.expected {
				t.Errorf("PERM.HasBasicUserPerm() = %v, want %v", got, tt.expected)
			}
		})
	}
}
