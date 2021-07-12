package bbs

import (
	"sync"
	"testing"
)

func TestUUserID_ToUsername(t *testing.T) {
	setupTest()
	defer teardownTest()

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

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := tt.u.ToUsername(); got != tt.expected {
				t.Errorf("UUserID.ToUsername() = %v, want %v", got, tt.expected)
			}
		})
	}

	wg.Wait()
}
