package bbs

import "testing"

func TestBBoardID_ToBrdname(t *testing.T) {
	tests := []struct {
		name     string
		b        BBoardID
		expected string
	}{
		// TODO: Add test cases.
		{
			b:        BBoardID("10_WhoAmI"),
			expected: "WhoAmI",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.ToBrdname(); got != tt.expected {
				t.Errorf("BBoardID.ToBrdname() = %v, want %v", got, tt.expected)
			}
		})
	}
}
