package types

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestIsDir(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		path string
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{"."},
			expected: true,
		},
		{
			args: args{"_non_exist_"},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := IsDir(tt.args.path); got != tt.expected {
				t.Errorf("IsDir() = %v, expected %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestCopyFileToFile(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{src: "./testcase/file1", dst: "./testcase/file"},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := CopyFileToFile(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("CopyFileToFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, err := os.Stat(tt.args.dst)
			if err != nil {
				t.Errorf("CopyFileToFile(): unable to find dst: src: %v dst: %v e: %v", tt.args.src, tt.args.dst, err)
			}

			os.RemoveAll(tt.args.dst)
		})
		wg.Wait()
	}
}

func TestCopyFile(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = MkdirAll("./testcase/to-dir")
	defer os.RemoveAll("./testcase/to-dir")

	_ = MkdirAll("./testcase/to-dir-2")
	defer os.RemoveAll("./testcase/to-dir-2")

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		src     string
		dst     string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "file-to-file",
			args: args{src: "./testcase/file1", dst: "./testcase/file"},
			src:  "./testcase/file1",
			dst:  "./testcase/file",
		},
		{
			name: "file-to-dir",
			args: args{src: "./testcase/file1", dst: "./testcase/to-dir"},
			src:  "./testcase/file1",
			dst:  "./testcase/to-dir/file1",
		},
		{
			name: "dir-to-dir",
			args: args{src: "./testcase/dir2", dst: "./testcase/to-dir-2"},
			src:  "./testcase/dir2/file2",
			dst:  "./testcase/to-dir-2/file2",
		},
		{
			name: "dir-to-dir (no dst-dir)",
			args: args{src: "./testcase/dir2", dst: "./testcase/to-dir-3"},
			src:  "./testcase/dir2/file2",
			dst:  "./testcase/to-dir-3/file2",
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := CopyFile(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			srcContent, _ := os.ReadFile(tt.src)

			_, err := os.Stat(tt.dst)
			if err != nil {
				t.Errorf("CopyFile() unable to find dst: src: %v dst: %v e: %v", tt.args.src, tt.args.dst, err)
				return
			}

			dstContent, _ := os.ReadFile(tt.dst)

			if !bytes.Equal(srcContent, dstContent) {
				t.Errorf("CopyFile() diff: src: %v dst: %v srcContent: %v dstContent: %v", tt.args.src, tt.args.dst, srcContent, dstContent)
				return
			}

			os.RemoveAll(tt.args.dst)
		})
		wg.Wait()
	}
}

func TestRename(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = MkdirAll("./testcase/to-dir")
	defer os.RemoveAll("./testcase/to-dir")

	_ = MkdirAll("./testcase/to-dir-2")
	defer os.RemoveAll("./testcase/to-dir-2")

	defer os.RemoveAll("./testcase/to-dir-3")
	defer os.RemoveAll("./testcase/to-dir-4")

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		src     string
		dst     string
		wantErr bool
	}{
		// TODO: Add test cases.
		// TODO: Add test cases.
		{
			name: "file-to-file",
			args: args{src: "./testcase/file1", dst: "./testcase/file"},
			src:  "./testcase/file1",
			dst:  "./testcase/file",
		},
		{
			name: "file-to-dir (force changing to-dir as file)",
			args: args{src: "./testcase/file1", dst: "./testcase/to-dir"},
			src:  "./testcase/file1",
			dst:  "./testcase/to-dir",
		},
		{
			name: "dir-to-dir",
			args: args{src: "./testcase/dir2", dst: "./testcase/to-dir-2"},
			src:  "./testcase/dir2/file2",
			dst:  "./testcase/to-dir-2/file2",
		},
		{
			name: "dir-to-dir (no dst-dir)",
			args: args{src: "./testcase/dir2", dst: "./testcase/to-dir-3"},
			src:  "./testcase/dir2/file2",
			dst:  "./testcase/to-dir-3/file2",
		},
		{
			name: "dir-to-dir (no dst-dir)",
			args: args{src: "./testcase/dir2", dst: "./testcase/to-dir-4/to-dir-5"},
			src:  "./testcase/dir2/file2",
			dst:  "./testcase/to-dir-4/to-dir-5/file2",
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			_ = CopyFile("./testcase/file1", "./testcase/file1.orig")
			defer Rename("./testcase/file1.orig", "./testcase/file1")

			_ = CopyFile("./testcase/dir2", "./testcase/dir2.orig")
			defer Rename("./testcase/dir2.orig", "testcase/dir2")

			srcContent, _ := os.ReadFile(tt.src)

			if err := Rename(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Rename() error = %v, wantErr %v", err, tt.wantErr)
			}

			dstContent, _ := os.ReadFile(tt.dst)

			if !bytes.Equal(srcContent, dstContent) {
				t.Errorf("Rename() diff: src: %v dst: %v srcContent: %v dstContent: %v", tt.src, tt.dst, srcContent, dstContent)
				return
			}

			os.RemoveAll(tt.args.dst)
		})
		wg.Wait()
	}
}

func TestMkdir(t *testing.T) {
	setupTest()
	defer teardownTest()

	defer os.RemoveAll("./testcase/test_dir0")
	defer os.RemoveAll("./testcase/test_dir1")

	path0 := "./testcase/test_dir0"
	path1 := "./testcase/test_dir1"

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{path: path0},
		},
		{
			args:    args{path: path0},
			wantErr: true,
		},
		{
			args: args{path: path1},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := Mkdir(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Mkdir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}

func TestBinaryRead(t *testing.T) {
	setupTest()
	defer teardownTest()

	type temp struct {
		A int32
		B [15]byte
		C int32
	}
	str0 := "\x01\x02\x00\x00abcdefghijklmno\x03\x04\x05\x00\x00"
	r0 := strings.NewReader(str0)
	data0 := &temp{}
	expectedB0 := [15]byte{}
	copy(expectedB0[:], []byte("abcdefghijklmno"))
	expected0 := &temp{
		A: 513,
		B: expectedB0,
		C: 1284,
	}

	str1 := "abcdefghijklmno"
	r1 := strings.NewReader(str1)
	data1 := [15]byte{}
	expected1 := [15]byte{}
	copy(expected1[:], expectedB0[:])

	str2 := "abcdefghijklmno"
	r2 := strings.NewReader(str2)
	data2 := make([]byte, 15)
	expected2 := make([]byte, 15)
	copy(expected2[:], expectedB0[:])

	type args struct {
		reader io.ReadSeeker
		order  binary.ByteOrder
		data   interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name:     "struct-ptr",
			args:     args{reader: r0, order: binary.LittleEndian, data: data0},
			expected: expected0,
		},
		{
			name:     "array",
			args:     args{reader: r1, order: binary.LittleEndian, data: &data1},
			expected: &expected1,
		},
		{
			name:     "slice",
			args:     args{reader: r2, order: binary.LittleEndian, data: data2},
			expected: expected2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BinaryRead(tt.args.reader, tt.args.order, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("BinaryRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
