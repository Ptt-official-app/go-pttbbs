package types

import (
	"reflect"
	"testing"
	"time"
)

func TestNowTS(t *testing.T) {
	start := Time4(time.Now().Unix())
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NowTS(); got < start {
				t.Errorf("NowTS() = %v, < start: %v", got, start)
			}
		})
	}
}

func TestTime4_ToLocal(t *testing.T) {
	tests := []struct {
		name     string
		tr       Time4
		expected time.Time
	}{
		// TODO: Add test cases.
		{
			tr:       Time4(1234567890),
			expected: time.Date(2009, 2, 14, 7, 31, 30, 0, TIMEZONE),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToLocal(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Time4.ToLocal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTime4_ToUtc(t *testing.T) {
	tests := []struct {
		name     string
		tr       Time4
		expected time.Time
	}{
		// TODO: Add test cases.
		{
			tr:       Time4(1234567890),
			expected: time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToUtc(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Time4.ToUtc() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTime4_Cdate(t *testing.T) {
	tests := []struct {
		name     string
		tr       Time4
		expected string
	}{
		// TODO: Add test cases.
		{
			tr:       Time4(1234567890),
			expected: "02/14/2009 07:31:30 Sat",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Cdate(); got != tt.expected {
				t.Errorf("Time4.Cdate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTime4_Cdatelite(t *testing.T) {
	tests := []struct {
		name     string
		tr       Time4
		expected string
	}{
		// TODO: Add test cases.
		{
			tr:       Time4(1234567890),
			expected: "02/14/2009 07:31:30",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Cdatelite(); got != tt.expected {
				t.Errorf("Time4.Cdatelite() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTime4_Cdatedate(t *testing.T) {
	tests := []struct {
		name     string
		tr       Time4
		expected string
	}{
		// TODO: Add test cases.
		{
			tr:       Time4(1234567890),
			expected: "02/14/2009",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Cdatedate(); got != tt.expected {
				t.Errorf("Time4.Cdatedate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTime4_CdateMd(t *testing.T) {
	tests := []struct {
		name     string
		tr       Time4
		expected string
	}{
		// TODO: Add test cases.
		{
			tr:       Time4(1234567890),
			expected: "02/14",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.CdateMd(); got != tt.expected {
				t.Errorf("Time4.CdateMd() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTime4_CdateMdHM(t *testing.T) {
	tests := []struct {
		name     string
		tr       Time4
		expected string
	}{
		// TODO: Add test cases.
		{
			tr:       Time4(1234567890),
			expected: "02/14 07:31",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.CdateMdHM(); got != tt.expected {
				t.Errorf("Time4.CdateMdHM() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTime4_CdateMdHMS(t *testing.T) {
	tests := []struct {
		name     string
		tr       Time4
		expected string
	}{
		// TODO: Add test cases.
		{
			tr:       Time4(1234567890),
			expected: "02/14 07:31:30",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.CdateMdHMS(); got != tt.expected {
				t.Errorf("Time4.CdateMdHMS() = %v, want %v", got, tt.expected)
			}
		})
	}
}
