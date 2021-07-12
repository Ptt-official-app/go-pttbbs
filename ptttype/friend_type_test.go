package ptttype

import (
	"sync"
	"testing"
)

func TestFriendType_Filename(t *testing.T) {
	tests := []struct {
		name     string
		f        FriendType
		expected string
	}{
		// TODO: Add test cases.
		{
			f:        FRIEND_OVERRIDE,
			expected: FN_OVERRIDES,
		},
		{
			f:        FRIEND_REJECT,
			expected: FN_REJECT,
		},
		{
			f:        FRIEND_ALOHA,
			expected: FN_ALOHAED,
		},
		{
			f:        FRIEND_SPECIAL,
			expected: "",
		},
		{
			f:        FRIEND_CANVOTE,
			expected: FN_CANVOTE,
		},
		{
			f:        BOARD_WATER,
			expected: FN_WATER,
		},
		{
			f:        BOARD_VISIBLE,
			expected: FN_VISIBLE,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := tt.f.Filename(); got != tt.expected {
				t.Errorf("FriendType.Filename() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}
