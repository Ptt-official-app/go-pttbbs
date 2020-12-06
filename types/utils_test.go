package types

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGetRandom(t *testing.T) {
	tests := []struct {
		name     string
		expected int
	}{
		// TODO: Add test cases.
		{expected: 22},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandom(); len(got) != tt.expected {
				t.Errorf("GetRandom() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetCurrentMilliTS(t *testing.T) {
	nowMilliTS := time.Now().UnixNano() / 1000000
	tests := []struct {
		name     string
		expected int64
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCurrentMilliTS(); got < nowMilliTS {
				t.Errorf("GetCurrentMilliTS() = %v, want %v", got, nowMilliTS)
			}
		})
	}
}

func TestIsRegularFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{"testcase/file1"},
			expected: true,
		},
		{
			args:     args{"testcase/non-exist.txt"},
			expected: false,
		},
		{
			args:     args{"testcase"},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRegularFile(tt.args.filename); got != tt.expected {
				t.Errorf("IsRegularFile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name     string
		args     args
		expected []byte
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{"testcase/file1"},
			expected: []byte("jalakjdhkasjdhfklsajdhf\x0a"),
		},
		{
			args:     args{"testcase/non-exist.txt"},
			expected: nil,
			wantErr:  true,
		},
		{
			args:     args{"testcase"},
			expected: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ReadFile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBinRead(t *testing.T) {
	theBytes := make([]byte, 3)
	theByte2 := make([]byte, 13)
	nextByte := [1]byte{}

	expect := []byte("012")

	type args struct {
		v       interface{}
		theSize uintptr
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantNext byte
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{theBytes, 6},
			expected: expect,
			wantNext: '6',
		},
		{
			name:    "bytes are more than the size",
			args:    args{theByte2, 3},
			wantErr: true,
		},
		{
			name:    "bytes are more than the size-2",
			args:    args{theBytes, 2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := os.Open("./testcase/bin_read_file.txt")
			defer f.Close()

			err := BinRead(f, tt.args.v, tt.args.theSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("BinRead() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(tt.expected, tt.args.v) {
				t.Errorf("BinRead() v: %v expected: %v", tt.args.v, tt.expected)
			}

			_, _ = f.Read(nextByte[:])

			if nextByte[0] != tt.wantNext {
				t.Errorf("BinRead() next: %v expected: %v", nextByte, tt.wantNext)
			}
		})
	}
}

func TestBinWrite(t *testing.T) {

	toWrite := [3]byte{}
	copy(toWrite[:], []byte("123"))
	toWrite2 := [5]byte{}
	copy(toWrite2[:], []byte("12345"))

	type args struct {
		v       interface{}
		theSize uintptr
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{&toWrite, 10},
		},
		{
			name:    "bytes to large",
			args:    args{&toWrite2, 3},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := os.OpenFile("testcase/bin_write_file.txt", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
			defer f.Close()

			err := BinWrite(f, tt.args.v, tt.args.theSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("BinWrite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			info, _ := os.Stat("testcase/bin_write_file.txt")
			nBytes := info.Size()
			if nBytes != int64(tt.args.theSize) {
				t.Errorf("BinWrite() nBytes: %v expected: %v", nBytes, tt.args.theSize)

			}
		})
	}
}

func TestReadLine(t *testing.T) {
	file, _ := os.Open("testcase/testreadline.txt")
	defer file.Close()

	file2, _ := os.Open("testcase/testreadline2.txt")
	defer file2.Close()

	reader := bufio.NewReader(file)
	reader2 := bufio.NewReader(file2)

	type args struct {
		reader *bufio.Reader
	}
	tests := []struct {
		name        string
		args        args
		expected    []byte
		wantErr     bool
		expectedErr error
	}{
		// TODO: Add test cases.
		{
			args:     args{reader},
			expected: []byte("test1"),
		},
		{
			args:     args{reader},
			expected: []byte("test2"),
		},
		{
			args:     args{reader},
			expected: []byte("test3"),
		},
		{
			args:        args{reader},
			expected:    nil,
			wantErr:     true,
			expectedErr: io.EOF,
		},
		{
			args:     args{reader2},
			expected: []byte("test1"),
		},
		{
			args:     args{reader2},
			expected: []byte("test2"),
		},
		{
			args:     args{reader2},
			expected: []byte("test3"),
		},
		{
			args:     args{reader2},
			expected: []byte("test4"),
		},
		{
			args:        args{reader},
			expected:    nil,
			wantErr:     true,
			expectedErr: io.EOF,
		},
		{
			args:        args{nil},
			expected:    nil,
			wantErr:     true,
			expectedErr: ErrNilReader,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadLine(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != tt.expectedErr {
				t.Errorf("ReadLine() e: %v expected: %v", err, tt.expectedErr)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ReadLine() = %v, want %v", got, tt.expected)
			}
		})
	}
}
