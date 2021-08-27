package ptttype

import (
	"testing"
)

func TestUFlag_HasUserFlag(t *testing.T) {
	type args struct {
		flag UFlag
	}
	tests := []struct {
		name     string
		u        UFlag
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			u:        UF_FRIEND | UF_ADBANNER,
			args:     args{UF_FRIEND},
			expected: true,
		},
		{
			u:        UF_ADBANNER,
			args:     args{UF_FRIEND},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.HasUserFlag(tt.args.flag); got != tt.expected {
				t.Errorf("UFlag.HasUserFlag() = %v, want %v", got, tt.expected)
			}
		})
	}
}
