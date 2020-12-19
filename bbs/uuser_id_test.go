package bbs

import "testing"

func TestUUserID_ToUsername(t *testing.T) {
	tests := []struct {
		name     string
		u        UUserID
		expected string
	}{
		// TODO: Add test cases.
		{
			u:        UUserID("SYSOP"),
			expected: "SYSOP",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.ToUsername(); got != tt.expected {
				t.Errorf("UUserID.ToUsername() = %v, want %v", got, tt.expected)
			}
		})
	}
}
