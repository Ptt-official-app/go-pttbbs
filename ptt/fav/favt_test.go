package fav

import (
	"testing"
)

func TestFavT_GetFav4TypeSize(t *testing.T) {
	tests := []struct {
		name     string
		theType  FavT
		expected uintptr
	}{
		// TODO: Add test cases.
		{
			theType:  FAVT_BOARD,
			expected: SIZE_OF_FAV4_BOARD,
		},
		{
			theType:  FAVT_FOLDER,
			expected: SIZE_OF_FAV_FOLDER,
		},
		{
			theType:  FAVT_LINE,
			expected: SIZE_OF_FAV_LINE,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.theType.GetFav4TypeSize(); got != tt.expected {
				t.Errorf("FavT.GetFav4TypeSize() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFavT_GetTypeSize(t *testing.T) {
	tests := []struct {
		name     string
		theType  FavT
		expected uintptr
	}{
		// TODO: Add test cases.
		{
			theType:  FAVT_BOARD,
			expected: SIZE_OF_FAV_BOARD,
		},
		{
			theType:  FAVT_FOLDER,
			expected: SIZE_OF_FAV_FOLDER,
		},
		{
			theType:  FAVT_LINE,
			expected: SIZE_OF_FAV_LINE,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.theType.GetTypeSize(); got != tt.expected {
				t.Errorf("FavT.GetTypeSize() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFavT_String(t *testing.T) {
	tests := []struct {
		name     string
		tr       FavT
		expected string
	}{
		// TODO: Add test cases.
		{
			tr:       FAVT_BOARD,
			expected: "Board",
		},
		{
			tr:       FAVT_FOLDER,
			expected: "Folder",
		},
		{
			tr:       FAVT_LINE,
			expected: "Line",
		},
		{
			tr:       0,
			expected: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.expected {
				t.Errorf("FavT.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
