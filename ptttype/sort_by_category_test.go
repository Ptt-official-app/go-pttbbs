package ptttype

import (
	"sync"
	"testing"
)

func TestSortByCategory_String(t *testing.T) {
	tests := []struct {
		name     string
		d        SortByCategory
		expected string
	}{
		// TODO: Add test cases.
		{
			d:        SORT_BY_ID,
			expected: "id",
		},
		{
			d:        SORT_BY_CLASS,
			expected: "class",
		},
		{
			d:        SORT_BY_STAT,
			expected: "stat",
		},
		{
			d:        SORT_BY_IDLE,
			expected: "idle",
		},
		{
			d:        SORT_BY_FROM,
			expected: "from",
		},
		{
			d:        SORT_BY_FIVE,
			expected: "five",
		},
		{
			d:        SORT_BY_SEX,
			expected: "gender",
		},
		{
			d:        SortByCategory(100),
			expected: "unknown",
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := tt.d.String(); got != tt.expected {
				t.Errorf("SortByCategory.String() = %v, want %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}
