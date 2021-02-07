package cmsys

import (
	"reflect"
	"testing"
)

func TestStringHashWithHashBits(t *testing.T) {
	type args struct {
		theBytes []byte
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312")},
			expected: 0x94b2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringHashWithHashBits(tt.args.theBytes); got != tt.expected {
				t.Errorf("StringHashWithHashBits() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringHash(t *testing.T) {
	type args struct {
		theBytes []byte
	}
	tests := []struct {
		name     string
		args     args
		expected uint32
	}{
		// TODO: Add test cases.
		{
			args:     args{[]byte("12312")},
			expected: 0x7aeb94b2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringHash(tt.args.theBytes); got != tt.expected {
				t.Errorf("StringHash() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStripNoneBig5(t *testing.T) {
	type args struct {
		str []byte
	}
	tests := []struct {
		name                 string
		args                 args
		expectedSanitizedStr []byte
	}{
		// TODO: Add test cases.
		{
			args:                 args{[]byte("12345")},
			expectedSanitizedStr: []byte("12345"),
		},
		{
			args:                 args{[]byte("\xff\xfd12345")},
			expectedSanitizedStr: []byte("\xff\xfd12345"),
		},
		{
			args:                 args{[]byte("\x80\x0112345")},
			expectedSanitizedStr: []byte("12345"),
		},
		{
			args:                 args{[]byte("\x80\x01\x80\x7312345")},
			expectedSanitizedStr: []byte("\x80\x7312345"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSanitizedStr := StripNoneBig5(tt.args.str); !reflect.DeepEqual(gotSanitizedStr, tt.expectedSanitizedStr) {
				t.Errorf("StripNoneBig5() = %v, want %v", gotSanitizedStr, tt.expectedSanitizedStr)
			}
		})
	}
}

func TestStripAnsi(t *testing.T) {
	src0 := []byte("asda\x1bqwem\x1b[1;37m12312312\x1b123121\x1b[m")
	expected0 := []byte("asdawem1231231223121")
	type args struct {
		src  []byte
		flag StripAnsiFlag
	}
	tests := []struct {
		name        string
		args        args
		expectedDst []byte
	}{
		// TODO: Add test cases.
		{
			args:        args{src: src0, flag: STRIP_ANSI_ALL},
			expectedDst: expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDst := StripAnsi(tt.args.src, tt.args.flag); !reflect.DeepEqual(gotDst, tt.expectedDst) {
				t.Errorf("StringAnsi() = %v, want %v", gotDst, tt.expectedDst)
			}
		})
	}
}

func TestStripBlank(t *testing.T) {
	type args struct {
		theBytes []byte
	}
	tests := []struct {
		name     string
		args     args
		expected []byte
	}{
		// TODO: Add test cases.
		{
			args:     args{theBytes: []byte("test")},
			expected: []byte("test"),
		},
		{
			args:     args{theBytes: []byte("test1 test2")},
			expected: []byte("test1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripBlank(tt.args.theBytes); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("StripBlank() = %v, want %v", got, tt.expected)
			}
		})
	}
}
