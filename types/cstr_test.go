package types

import (
	"reflect"
	"testing"
)

func TestCstrToBytes(t *testing.T) {
	setupTest()
	defer teardownTest()

	str1 := [13]byte{}
	str2 := [13]byte{}
	copy(str2[:], []byte("123"))
	str3 := [10]byte{}
	copy(str3[:], []byte("0123456789"))
	str4 := [10]byte{}
	copy(str4[:], []byte("01234\x006789"))

	type args struct {
		cstr Cstr
	}
	tests := []struct {
		name     string
		args     args
		expected []byte
	}{
		{
			name:     "init",
			args:     args{str1[:]},
			expected: []byte{},
		},
		{
			name:     "with only 3 letters",
			args:     args{str2[:]},
			expected: []byte("123"),
		},
		{
			name:     "with no 0",
			args:     args{str3[:]},
			expected: []byte("0123456789"),
		},
		{
			name:     "cutoff at str4[5]",
			args:     args{str4[:]},
			expected: []byte("01234"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CstrToBytes(tt.args.cstr); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("CstrToBytes() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCstrToString(t *testing.T) {
	setupTest()
	defer teardownTest()

	str1 := [13]byte{}
	str2 := [13]byte{}
	copy(str2[:], []byte("123"))
	str3 := [10]byte{}
	copy(str3[:], []byte("0123456789"))
	str4 := [10]byte{}
	copy(str4[:], []byte("01234\x006789"))

	type args struct {
		cstr Cstr
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		// TODO: Add test cases.
		{
			name:     "init",
			args:     args{str1[:]},
			expected: "",
		},
		{
			name:     "with only 3 letters",
			args:     args{str2[:]},
			expected: "123",
		},
		{
			name:     "with no 0",
			args:     args{str3[:]},
			expected: "0123456789",
		},
		{
			name:     "cutoff at str4[5]",
			args:     args{str4[:]},
			expected: "01234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CstrToString(tt.args.cstr); got != tt.expected {
				t.Errorf("CstrToString() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCstrcmp(t *testing.T) {
	type args struct {
		cstr1 Cstr
		cstr2 Cstr
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		// TODO: Add test cases.
		{
			name:     "eq",
			args:     args{Cstr([]byte("12345")), Cstr([]byte("12345"))},
			expected: 0,
		},
		{
			name:     "gt with eq length",
			args:     args{Cstr([]byte("34567")), Cstr([]byte("12345"))},
			expected: 2,
		},
		{
			name:     "gt with diff-length",
			args:     args{Cstr([]byte("34567")), Cstr([]byte("123456"))},
			expected: 2,
		},
		{
			name:     "lt with eq length",
			args:     args{Cstr([]byte("12345")), Cstr([]byte("34567"))},
			expected: -2,
		},
		{
			name:     "lt with diff length",
			args:     args{Cstr([]byte("1234567")), Cstr([]byte("34567"))},
			expected: -2,
		},
		{
			name:     "eq with \x00",
			args:     args{Cstr([]byte("123\x0056")), Cstr([]byte("123\x0045"))},
			expected: 0,
		},
		{
			name:     "eq with \x00",
			args:     args{Cstr([]byte("123\x0056")), Cstr([]byte("123\x004"))},
			expected: 0,
		},
		{
			name:     "eq with \x00",
			args:     args{Cstr([]byte("123\x005")), Cstr([]byte("123\x0056"))},
			expected: 0,
		},
		{
			name:     "len(c1) > len(c2)",
			args:     args{Cstr([]byte("12345")), Cstr([]byte("123"))},
			expected: 52,
		},
		{
			name:     "len(c1) > len(c2) with \x00",
			args:     args{Cstr([]byte("12345")), Cstr([]byte("123\x0045"))},
			expected: 52,
		},
		{
			name:     "len(c1) > len(c2) with \x00",
			args:     args{Cstr([]byte("1234\x005")), Cstr([]byte("123\x0045"))},
			expected: 52,
		},
		{
			name:     "len(c1) > len(c2)",
			args:     args{Cstr([]byte("12345")), Cstr([]byte("123"))},
			expected: 52,
		},
		{
			name:     "len(c1) > len(c2) with \x00 ('4' - 0)",
			args:     args{Cstr([]byte("12345")), Cstr([]byte("123\x0045"))},
			expected: 52,
		},
		{
			name:     "len(c1) > len(c2) with \x00 ('4' - 0)",
			args:     args{Cstr([]byte("1234\x005")), Cstr([]byte("123\x0045"))},
			expected: 52,
		},
		{
			name:     "len(c1) < len(c2)",
			args:     args{Cstr([]byte("123")), Cstr([]byte("1234"))},
			expected: -52,
		},
		{
			name:     "len(c1) < len(c2) with \x00",
			args:     args{Cstr([]byte("123\x0045")), Cstr([]byte("12345"))},
			expected: -52,
		},
		{
			name:     "len(c1) < len(c2) with \x00",
			args:     args{Cstr([]byte("123\x0045")), Cstr([]byte("1234\x005"))},
			expected: -52,
		},
		{
			name:     "len(c1) < len(c2)",
			args:     args{Cstr([]byte("123")), Cstr([]byte("12345"))},
			expected: -52,
		},
		{
			name:     "len(c1) < len(c2) with \x00 (0 - '4')",
			args:     args{Cstr([]byte("123\x0045")), Cstr([]byte("12345"))},
			expected: -52,
		},
		{
			name:     "len(c1) < len(c2) with \x00",
			args:     args{Cstr([]byte("123\x0045")), Cstr([]byte("1234\x005"))},
			expected: -52,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cstrcmp(tt.args.cstr1, tt.args.cstr2); got != tt.expected {
				t.Errorf("Cstrcmp() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCstrTolower(t *testing.T) {
	type args struct {
		cstr Cstr
	}
	tests := []struct {
		name     string
		args     args
		expected Cstr
	}{
		// TODO: Add test cases.
		{
			args:     args{Cstr([]byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*()_-+=[]{}"))},
			expected: Cstr([]byte("0123456789abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz!@#$%^&*()_-+=[]{}")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CstrTolower(tt.args.cstr); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("CstrTolower() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCstrstr(t *testing.T) {
	type args struct {
		cstr   Cstr
		substr Cstr
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		// TODO: Add test cases.
		{
			args:     args{Cstr([]byte("123\x0045")), Cstr([]byte("34"))},
			expected: -1,
		},
		{
			args:     args{Cstr([]byte("12345")), Cstr([]byte("34"))},
			expected: 2,
		},
		{
			args:     args{Cstr([]byte("12345")), Cstr([]byte("12"))},
			expected: 0,
		},
		{
			args:     args{Cstr([]byte("12345")), Cstr([]byte("45"))},
			expected: 3,
		},
		{
			args:     args{Cstr([]byte("12345")), Cstr([]byte("012"))},
			expected: -1,
		},
		{
			args:     args{Cstr([]byte("12345")), Cstr([]byte("456"))},
			expected: -1,
		},
		{
			args: args{
				cstr:   Cstr([]byte("456")),
				substr: Cstr([]byte("4567")),
			},
			expected: -1,
		},
		{
			args: args{
				cstr:   Cstr([]byte("bc")),
				substr: Cstr([]byte("abc")),
			},
			expected: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cstrstr(tt.args.cstr, tt.args.substr); got != tt.expected {
				t.Errorf("Cstrstr() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCstrcasecmp(t *testing.T) {
	type args struct {
		cstr1 Cstr
		cstr2 Cstr
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		// TODO: Add test cases.
		{
			args:     args{Cstr("abc"), Cstr("Abc")},
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cstrcasecmp(tt.args.cstr1, tt.args.cstr2); got != tt.expected {
				t.Errorf("Cstrcasecmp() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCstrcasestr(t *testing.T) {
	type args struct {
		cstr   Cstr
		substr Cstr
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		// TODO: Add test cases.
		{
			args:     args{Cstr("abc"), Cstr("Abc")},
			expected: 0,
		},
		{
			args: args{
				cstr:   Cstr("bc"),
				substr: Cstr("Abc"),
			},
			expected: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cstrcasestr(tt.args.cstr, tt.args.substr); got != tt.expected {
				t.Errorf("Cstrcasestr() = %v, want %v cstr: %v substr: %v", got, tt.expected, tt.args.cstr, tt.args.substr)
			}
		})
	}
}
