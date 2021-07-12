package bbs

import (
	"sync"
	"testing"
)

func TestBBoardID_ToBrdname(t *testing.T) {
	setupTest()
	defer teardownTest()

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
		{
			b:        BBoardID("50_C_Chat"),
			expected: "C_Chat",
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := tt.b.ToBrdname(); got != tt.expected {
				t.Errorf("BBoardID.ToBrdname() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}
