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

func TestStrcaseStartsWith(t *testing.T) {
	str0 := []byte{'t', 'e', 's', 't'}
	prefix0 := []byte{'T'}
	expected0 := true

	str1 := []byte{'t', 'e', 's', 't'}
	prefix1 := []byte{'t', 'E'}
	expected1 := true

	str2 := []byte{'t', 'e', 's', 't'}
	prefix2 := []byte{'t', 'e'}
	expected2 := true

	str3 := []byte{'t', 'e', 's', 't'}
	prefix3 := []byte{'t', 'e', 's', 't', '3'}
	expected3 := false

	type args struct {
		str    []byte
		prefix []byte
	}
	tests := []struct {
		name            string
		args            args
		expectedIsValid bool
	}{
		// TODO: Add test cases.
		{
			args:            args{str: str0, prefix: prefix0},
			expectedIsValid: expected0,
		},
		{
			args:            args{str: str1, prefix: prefix1},
			expectedIsValid: expected1,
		},
		{
			args:            args{str: str2, prefix: prefix2},
			expectedIsValid: expected2,
		},
		{
			args:            args{str: str3, prefix: prefix3},
			expectedIsValid: expected3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsValid := StrcaseStartsWith(tt.args.str, tt.args.prefix); gotIsValid != tt.expectedIsValid {
				t.Errorf("StrcaseStartsWith() = %v, want %v", gotIsValid, tt.expectedIsValid)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	type args struct {
		str []byte
	}
	tests := []struct {
		name           string
		args           args
		expectedNewStr []byte
	}{
		// TODO: Add test cases.
		{
			args:           args{str: []byte{'a', 'b', 'c', ' ', ' ', ' ', ' '}},
			expectedNewStr: []byte{'a', 'b', 'c'},
		},
		{
			args:           args{str: []byte{'a', 'b', 'c', ' ', ' ', ' ', ' ', 0, ' ', ' ', ' '}},
			expectedNewStr: []byte{'a', 'b', 'c'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNewStr := Trim(tt.args.str); !reflect.DeepEqual(gotNewStr, tt.expectedNewStr) {
				t.Errorf("Trim() = %v, want %v", gotNewStr, tt.expectedNewStr)
			}
		})
	}
}

func TestDBCSSafeTrim(t *testing.T) {
	type args struct {
		str []byte
	}
	tests := []struct {
		name           string
		args           args
		expectedNewStr []byte
	}{
		// TODO: Add test cases.
		{
			args:           args{str: []byte{0xbc, 0xd0, 0xc3, 0x44, 0x3a}},
			expectedNewStr: []byte{0xbc, 0xd0, 0xc3, 0x44, 0x3a},
		},
		{
			args:           args{str: []byte{0xbc, 0xd0, 0xc3}},
			expectedNewStr: []byte{0xbc, 0xd0},
		},
		{
			args:           args{str: []byte{0xbc, 0xd0, 0xc3, 0x44}},
			expectedNewStr: []byte{0xbc, 0xd0, 0xc3, 0x44},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNewStr := DBCSSafeTrim(tt.args.str); !reflect.DeepEqual(gotNewStr, tt.expectedNewStr) {
				t.Errorf("DBCSSafeTrim() = %v, want %v", gotNewStr, tt.expectedNewStr)
			}
		})
	}
}

func TestDBCSStatus(t *testing.T) {
	type args struct {
		str []byte
		pos int
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus DBCSStatus_t
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStatus := DBCSStatus(tt.args.str, tt.args.pos); gotStatus != tt.expectedStatus {
				t.Errorf("DBCSStatus() = %v, want %v", gotStatus, tt.expectedStatus)
			}
		})
	}
}

func TestDBCSNextStatus(t *testing.T) {
	type args struct {
		c          byte
		prevStatus DBCSStatus_t
	}
	tests := []struct {
		name              string
		args              args
		expectedNewStatus DBCSStatus_t
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNewStatus := DBCSNextStatus(tt.args.c, tt.args.prevStatus); gotNewStatus != tt.expectedNewStatus {
				t.Errorf("DBCSNextStatus() = %v, want %v", gotNewStatus, tt.expectedNewStatus)
			}
		})
	}
}
