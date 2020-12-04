package ptttype

import "testing"

func TestMsgMode_String(t *testing.T) {
	tests := []struct {
		name     string
		m        MsgMode
		expected string
	}{
		// TODO: Add test cases.
		{
			m:        MSGMODE_TALK,
			expected: "talk",
		},
		{
			m:        MSGMODE_WRITE,
			expected: "write",
		},
		{
			m:        MSGMODE_FROMANGEL,
			expected: "from-angel",
		},
		{
			m:        MSGMODE_TOANGEL,
			expected: "to-angel",
		},
		{
			m:        MsgMode(100),
			expected: "[unknown]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.expected {
				t.Errorf("MsgMode.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
