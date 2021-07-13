package cmsys

import (
	"encoding/binary"
	"io"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestGetNumRecords(t *testing.T) {
	type args struct {
		filename string
		size     uintptr
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		// TODO: Add test cases.
		{
			args:     args{filename: "testcase/BOARD", size: ptttype.BOARD_HEADER_RAW_SZ},
			expected: 12,
		},
		{
			args:     args{filename: "testcase/not-exist", size: ptttype.BOARD_HEADER_RAW_SZ},
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNumRecords(tt.args.filename, tt.args.size); got != tt.expected {
				t.Errorf("GetNumRecords() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetRecords(t *testing.T) {
	boardID := &ptttype.BoardID_t{}
	copy(boardID[:], []byte("WhoAmI"))
	filename0 := "./testcase/DIR"
	n0 := 100
	expected0 := []*ptttype.ArticleSummaryRaw{testArticleSummary0}

	filename1 := "./testcase/DIR1"
	startIdx1 := ptttype.SortIdx(1)
	n1 := 100
	expected1 := []*ptttype.ArticleSummaryRaw{testArticleSummary1, testArticleSummary2}

	filename2 := "./testcase/DIR"
	startIdx2 := ptttype.SortIdx(2)
	n2 := 100
	expected2 := []*ptttype.ArticleSummaryRaw{}

	filename3 := "./testcase/DIR"
	startIdx3 := ptttype.SortIdx(3)
	n3 := 100
	expected3 := []*ptttype.ArticleSummaryRaw{}

	fileHeaders := []ptttype.FileHeaderRaw{
		*testArticleSummary1.FileHeaderRaw,
		*testArticleSummary2.FileHeaderRaw,
		*testArticleSummary3.FileHeaderRaw,
		*testArticleSummary4.FileHeaderRaw,
		*testArticleSummary5.FileHeaderRaw,
	}

	filename4 := "./testcase/DIR_GET_RECORDS"
	defer os.RemoveAll(filename4)
	file, _ := os.OpenFile(filename4, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	_ = types.BinaryWrite(file, binary.LittleEndian, fileHeaders)

	startIdx4 := ptttype.SortIdx(2)
	n4 := 100
	expected4 := []*ptttype.ArticleSummaryRaw{
		testArticleSummary2,
		testArticleSummary3,
		testArticleSummary4,
		testArticleSummary5,
	}

	startIdx5 := ptttype.SortIdx(5)
	n5 := 100
	expected5 := []*ptttype.ArticleSummaryRaw{
		testArticleSummary5,
		testArticleSummary4,
		testArticleSummary3,
		testArticleSummary2,
		testArticleSummary1,
	}
	isDesc5 := true

	type args struct {
		boardID  *ptttype.BoardID_t
		filename string
		startIdx ptttype.SortIdx
		n        int
		isDesc   bool
	}
	tests := []struct {
		name     string
		args     args
		expected []*ptttype.ArticleSummaryRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{boardID: boardID, filename: filename0, startIdx: 1, n: n0, isDesc: false},
			expected: expected0,
		},
		{
			args:     args{boardID: boardID, filename: filename1, startIdx: startIdx1, n: n1, isDesc: false},
			expected: expected1,
		},
		{
			args:     args{boardID: boardID, filename: filename2, startIdx: startIdx2, n: n2, isDesc: false},
			expected: expected2,
		},
		{
			args:     args{boardID: boardID, filename: filename3, startIdx: startIdx3, n: n3, isDesc: false},
			expected: expected3,
		},
		{
			args:     args{boardID: boardID, filename: filename4, startIdx: startIdx4, n: n4, isDesc: false},
			expected: expected4,
		},
		{
			args:     args{boardID: boardID, filename: filename4, startIdx: startIdx5, n: n5, isDesc: isDesc5},
			expected: expected5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRecords(tt.args.boardID, tt.args.filename, tt.args.startIdx, tt.args.n, tt.args.isDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, each := range got {
				each.BoardID = nil
			}
			testutil.TDeepEqual(t, "fileheader", got, tt.expected)
		})
	}
}

func TestFindRecordStartAid(t *testing.T) {
	fileHeaders := []ptttype.FileHeaderRaw{
		*testArticleSummary1.FileHeaderRaw,
		*testArticleSummary2.FileHeaderRaw,
		*testArticleSummary3.FileHeaderRaw,
		*testArticleSummary4.FileHeaderRaw,
		*testArticleSummary5.FileHeaderRaw,
	}

	filename := "./testcase/DIR_FIND_RECORD_START_IDX"
	defer os.RemoveAll(filename)
	file, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	_ = types.BinaryWrite(file, binary.LittleEndian, fileHeaders)

	type args struct {
		dirFilename string
		total       int
		createTime  types.Time4
		filename    *ptttype.Filename_t
		isDesc      bool
	}
	tests := []struct {
		name             string
		args             args
		expectedStartIdx ptttype.SortIdx
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			name:             "find data before top with isDesc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607202239, filename: nil, isDesc: true},
			expectedStartIdx: -1,
		},
		{
			name:             "find data with same CreateTime with no isDesc, but with diff filename",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203395, filename: nil, isDesc: false},
			expectedStartIdx: 5,
		},
		{
			name:             "find data with same CreateTime with isDesc, but with diff filename",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203395, filename: nil, isDesc: true},
			expectedStartIdx: 1,
		},
		{
			name:             "find data after bottom with no isDesc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203396, filename: nil, isDesc: false},
			expectedStartIdx: -1,
		},
		{
			name:             "1st-desc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607202239, filename: &testArticleSummary1.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 1,
		},
		{
			name:             "2nd-desc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203395, filename: &testArticleSummary2.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 2,
		},
		{
			name:             "3rd-desc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203395, filename: &testArticleSummary3.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "4th-desc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203395, filename: &testArticleSummary4.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 4,
		},
		{
			name:             "5th-desc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203396, filename: &testArticleSummary5.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 5,
		},
		{
			name:             "1st-asc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607202239, filename: &testArticleSummary1.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 1,
		},
		{
			name:             "2nd-asc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203395, filename: &testArticleSummary2.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 2,
		},
		{
			name:             "3rd-asc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203395, filename: &testArticleSummary3.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "4th-asc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203395, filename: &testArticleSummary4.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 4,
		},
		{
			name:             "5th-asc",
			args:             args{dirFilename: filename, total: 5, createTime: 1607203396, filename: &testArticleSummary5.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStartIdx, err := FindRecordStartIdx(tt.args.dirFilename, tt.args.total, tt.args.createTime, tt.args.filename, tt.args.isDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindRecordStartIdx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStartIdx, tt.expectedStartIdx) {
				t.Errorf("FindRecordStartIdx() = %v, want %v", int(gotStartIdx), int(tt.expectedStartIdx))
			}
		})
	}
}

func TestSubstituteRecord(t *testing.T) {
	filename := "./testcase/testSubstituteRecord.txt"
	defer os.Remove(filename)

	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()

	types.BinaryWrite(file, binary.LittleEndian, testArticleSummary0.FileHeaderRaw)
	types.BinaryWrite(file, binary.LittleEndian, testArticleSummary1.FileHeaderRaw)
	types.BinaryWrite(file, binary.LittleEndian, testArticleSummary2.FileHeaderRaw)

	type args struct {
		filename   string
		data       interface{}
		theSize    uintptr
		idxInStore int32
	}
	tests := []struct {
		name     string
		args     args
		expected *ptttype.FileHeaderRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{filename: filename, data: testArticleSummary3.FileHeaderRaw, theSize: ptttype.FILE_HEADER_RAW_SZ, idxInStore: 1},
			expected: testArticleSummary3.FileHeaderRaw,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := SubstituteRecord(tt.args.filename, tt.args.data, tt.args.theSize, tt.args.idxInStore); (err != nil) != tt.wantErr {
				t.Errorf("SubstituteRecord() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := &ptttype.FileHeaderRaw{}
			offset := int64(ptttype.FILE_HEADER_RAW_SZ) * int64(tt.args.idxInStore)
			_, _ = file.Seek(offset, io.SeekStart)
			_ = types.BinaryRead(file, binary.LittleEndian, got)
			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
		wg.Wait()
	}
}

func TestAppendRecord(t *testing.T) {
	filename := "./testcase/testAppendRecord.txt"
	defer os.Remove(filename)

	type args struct {
		filename string
		data     interface{}
		theSize  uintptr
	}
	tests := []struct {
		name        string
		args        args
		expectedIdx ptttype.SortIdx
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{filename: filename, data: testArticleSummary0.FileHeaderRaw, theSize: ptttype.FILE_HEADER_RAW_SZ},
			expectedIdx: 1,
		},
		{
			args:        args{filename: filename, data: testArticleSummary1.FileHeaderRaw, theSize: ptttype.FILE_HEADER_RAW_SZ},
			expectedIdx: 2,
		},
		{
			args:        args{filename: filename, data: testArticleSummary2.FileHeaderRaw, theSize: ptttype.FILE_HEADER_RAW_SZ},
			expectedIdx: 3,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotIdx, err := AppendRecord(tt.args.filename, tt.args.data, tt.args.theSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppendRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIdx, tt.expectedIdx) {
				t.Errorf("AppendRecord() = %v, want %v", gotIdx, tt.expectedIdx)
			}
		})
		wg.Wait()
	}
}
