package ptttype

import (
	"sync"
	"testing"
)

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

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := tt.m.String(); got != tt.expected {
				t.Errorf("MsgMode.String() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}
