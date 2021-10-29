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
	file, _ := os.OpenFile(filename4, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
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

func TestFindRecordStartIdx(t *testing.T) {
	fileHeaders0 := []ptttype.FileHeaderRaw{
		*testArticleSummary1.FileHeaderRaw, // M.1607202239.A.30D
		*testArticleSummary2.FileHeaderRaw, // M.1607203395.A.F6C
		*testArticleSummary3.FileHeaderRaw, // M.1607203395.A.F6D
		*testArticleSummary4.FileHeaderRaw, // M.1607203395.A.F6A
		*testArticleSummary5.FileHeaderRaw, // M.1607203396.A.F6A
	}
	filename0 := "./testcase/DIR_FIND_RECORD_START_IDX0"
	defer os.RemoveAll(filename0)
	file0, _ := os.OpenFile(filename0, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	defer file0.Close()
	_ = types.BinaryWrite(file0, binary.LittleEndian, fileHeaders0)

	fileHeaders1 := []ptttype.FileHeaderRaw{
		*testArticleSummary1.FileHeaderRaw,   // M.1607202239.A.30D
		*testArticleSummary2_d.FileHeaderRaw, // .d1607203395.A.F6C
		*testArticleSummary3_d.FileHeaderRaw, // .d1607203395.A.F6D
		*testArticleSummary4.FileHeaderRaw,   // M.1607203395.A.F6A
		*testArticleSummary5.FileHeaderRaw,   // M.1607203396.A.F6A
	}
	filename1 := "./testcase/DIR_FIND_RECORD_START_IDX1"
	defer os.RemoveAll(filename1)
	file1, _ := os.OpenFile(filename1, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	defer file1.Close()
	_ = types.BinaryWrite(file1, binary.LittleEndian, fileHeaders1)

	fileHeaders2 := []ptttype.FileHeaderRaw{
		*testArticleSummary1.FileHeaderRaw,    // M.1607202239.A.30D
		*testArticleSummary2_d2.FileHeaderRaw, // .deleted
		*testArticleSummary3_d.FileHeaderRaw,  // .d1607203395.A.F6D
		*testArticleSummary4_d2.FileHeaderRaw, // .deleted
		*testArticleSummary5.FileHeaderRaw,    // M.1607203396.A.F6A
	}
	filename2 := "./testcase/DIR_FIND_RECORD_START_IDX2"
	defer os.RemoveAll(filename2)
	file2, _ := os.OpenFile(filename2, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	defer file2.Close()
	_ = types.BinaryWrite(file2, binary.LittleEndian, fileHeaders2)

	searchFilename3 := &ptttype.Filename_t{}
	copy(searchFilename3[:], []byte("M.1607203395.A.F6F"))

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
			name:             "find data before top with isDesc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607202238, filename: nil, isDesc: true},
			expectedStartIdx: -1,
			wantErr:          true,
		},
		{
			name:             "find data with same CreateTime with isAsc (not isDesc), but with nil filename (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: nil, isDesc: false},
			expectedStartIdx: 2,
		},
		{
			name:             "find data with same CreateTime, boundary case (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607202239, filename: nil, isDesc: false},
			expectedStartIdx: 1,
		},
		{
			name:             "find data with same CreateTime, boundary case (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607202239, filename: nil, isDesc: true},
			expectedStartIdx: 1,
		},
		{
			name:             "find data with same CreateTime, boundary case (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203396, filename: nil, isDesc: false},
			expectedStartIdx: 5,
		},
		{
			name:             "find data with same CreateTime, boundary case (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203396, filename: nil, isDesc: true},
			expectedStartIdx: 5,
		},
		{
			name:             "find data with same CreateTime with isDesc, but with nil filename (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: nil, isDesc: true},
			expectedStartIdx: 4,
		},
		{
			name:             "find data after bottom with isAsc (not isDesc) (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203397, filename: nil, isDesc: false},
			expectedStartIdx: -1,
			wantErr:          true,
		},
		{
			name:             "1st-desc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607202239, filename: &testArticleSummary1.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 1,
		},
		{
			name:             "2nd-desc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: &testArticleSummary2.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 2,
		},
		{
			name:             "3rd-desc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: &testArticleSummary3.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "4th-desc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: &testArticleSummary4.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 4,
		},
		{
			name:             "5th-desc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203396, filename: &testArticleSummary5.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 5,
		},
		{
			name:             "1st-asc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607202239, filename: &testArticleSummary1.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 1,
		},
		{
			name:             "2nd-asc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: &testArticleSummary2.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 2,
		},
		{
			name:             "3rd-asc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: &testArticleSummary3.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 3,
		},
		{
			name:             "4th-asc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: &testArticleSummary4.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 4,
		},
		{
			name:             "5th-asc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203396, filename: &testArticleSummary5.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 5,
		},
		{
			name:             "find filename3 with isDesc (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: searchFilename3, isDesc: true},
			expectedStartIdx: 4,
		},
		{
			name:             "find filename3 with asc (not isDesc) (file0)",
			args:             args{dirFilename: filename0, total: 5, createTime: 1607203395, filename: searchFilename3, isDesc: false},
			expectedStartIdx: 2,
		},

		// filename1
		{
			name:             "find data before top with isDesc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607202238, filename: nil, isDesc: true},
			expectedStartIdx: -1,
			wantErr:          true,
		},
		{
			name:             "find data with same CreateTime with isAsc (not isDesc), but with nil filename (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: nil, isDesc: false},
			expectedStartIdx: 2,
		},
		{
			name:             "find data with same CreateTime, boundary case (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607202239, filename: nil, isDesc: false},
			expectedStartIdx: 1,
		},
		{
			name:             "find data with same CreateTime, boundary case (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607202239, filename: nil, isDesc: true},
			expectedStartIdx: 1,
		},
		{
			name:             "find data with same CreateTime, boundary case (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203396, filename: nil, isDesc: false},
			expectedStartIdx: 5,
		},
		{
			name:             "find data with same CreateTime, boundary case (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203396, filename: nil, isDesc: true},
			expectedStartIdx: 5,
		},
		{
			name:             "find data with same CreateTime with isDesc, but with nil filename (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: nil, isDesc: true},
			expectedStartIdx: 4,
		},
		{
			name:             "find data after bottom with isAsc (not isDesc) (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203397, filename: nil, isDesc: false},
			expectedStartIdx: -1,
			wantErr:          true,
		},
		{
			name:             "1st-desc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607202239, filename: &testArticleSummary1.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 1,
		},
		{
			name:             "2nd-desc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: &testArticleSummary2.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 2,
		},
		{
			name:             "3rd-desc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: &testArticleSummary3.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "4th-desc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: &testArticleSummary4.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 4,
		},
		{
			name:             "5th-desc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203396, filename: &testArticleSummary5.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 5,
		},
		{
			name:             "1st-asc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607202239, filename: &testArticleSummary1.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 1,
		},
		{
			name:             "2nd-asc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: &testArticleSummary2.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 2,
		},
		{
			name:             "3rd-asc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: &testArticleSummary3.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 3,
		},
		{
			name:             "4th-asc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: &testArticleSummary4.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 4,
		},
		{
			name:             "5th-asc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203396, filename: &testArticleSummary5.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 5,
		},
		{
			name:             "find filename3 with isDesc (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: searchFilename3, isDesc: true},
			expectedStartIdx: 4,
		},
		{
			name:             "find filename3 with asc (not isDesc) (file1)",
			args:             args{dirFilename: filename1, total: 5, createTime: 1607203395, filename: searchFilename3, isDesc: false},
			expectedStartIdx: 2,
		},

		// filename2
		{
			name:             "find data before top with isDesc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607202238, filename: nil, isDesc: true},
			expectedStartIdx: -1,
			wantErr:          true,
		},
		{
			name:             "find data with same CreateTime with isAsc (not isDesc), but with nil filename (2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: nil, isDesc: false},
			expectedStartIdx: 3,
		},
		{
			name:             "find data with same CreateTime, boundary case (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607202239, filename: nil, isDesc: false},
			expectedStartIdx: 1,
		},
		{
			name:             "find data with same CreateTime, boundary case (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607202239, filename: nil, isDesc: true},
			expectedStartIdx: 1,
		},
		{
			name:             "find data with same CreateTime, boundary case (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203396, filename: nil, isDesc: false},
			expectedStartIdx: 5,
		},
		{
			name:             "find data with same CreateTime, boundary case (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203396, filename: nil, isDesc: true},
			expectedStartIdx: 5,
		},
		{
			name:             "find data with same CreateTime with isDesc, but with nil filename (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: nil, isDesc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "find data after bottom with isAsc (not isDesc) (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203397, filename: nil, isDesc: false},
			expectedStartIdx: -1,
			wantErr:          true,
		},
		{
			name:             "1st-desc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607202239, filename: &testArticleSummary1.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 1,
		},
		{
			name:             "2nd-desc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: &testArticleSummary2.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "3rd-desc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: &testArticleSummary3.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "4th-desc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: &testArticleSummary4.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "5th-desc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203396, filename: &testArticleSummary5.FileHeaderRaw.Filename, isDesc: true},
			expectedStartIdx: 5,
		},
		{
			name:             "1st-asc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607202239, filename: &testArticleSummary1.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 1,
		},
		{
			name:             "2nd-asc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: &testArticleSummary2.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 3,
		},
		{
			name:             "3rd-asc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: &testArticleSummary3.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 3,
		},
		{
			name:             "4th-asc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: &testArticleSummary4.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 3,
		},
		{
			name:             "5th-asc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203396, filename: &testArticleSummary5.FileHeaderRaw.Filename, isDesc: false},
			expectedStartIdx: 5,
		},

		{
			name:             "find filename3 with isDesc (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: searchFilename3, isDesc: true},
			expectedStartIdx: 3,
		},
		{
			name:             "find filename3 with asc (not isDesc) (file2)",
			args:             args{dirFilename: filename2, total: 5, createTime: 1607203395, filename: searchFilename3, isDesc: false},
			expectedStartIdx: 3,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotStartIdx, err := FindRecordStartIdx(tt.args.dirFilename, tt.args.total, tt.args.createTime, tt.args.filename, tt.args.isDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindRecordStartIdx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStartIdx, tt.expectedStartIdx) {
				t.Errorf("FindRecordStartIdx() = %v, want %v", int(gotStartIdx), int(tt.expectedStartIdx))
			}
		})
		wg.Wait()
	}
}

func TestSubstituteRecord(t *testing.T) {
	filename := "./testcase/testSubstituteRecord.txt"
	defer os.Remove(filename)

	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
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

func TestDeleteRecord(t *testing.T) {
	case_1_FileHeaders := []ptttype.FileHeaderRaw{
		*testArticleSummary1.FileHeaderRaw, // M.1607202239.A.30D
		*testArticleSummary2.FileHeaderRaw, // M.1607203395.A.F6C
		*testArticleSummary3.FileHeaderRaw, // M.1607203395.A.F6D
		*testArticleSummary4.FileHeaderRaw, // M.1607203395.A.F6A
		*testArticleSummary5.FileHeaderRaw, // M.1607203396.A.F6A
	}
	case_1_Filename := "./testcase/DIR_DELETE_ARTICLE"
	defer os.RemoveAll(case_1_Filename)
	file, _ := os.OpenFile(case_1_Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	defer file.Close()
	_ = types.BinaryWrite(file, binary.LittleEndian, case_1_FileHeaders)

	type args struct {
		filename string
		index    ptttype.SortIdxInStore
		theSize  uintptr
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test delete index 0 w/o error",
			args: args{
				filename: case_1_Filename,
				index:    0,
				theSize:  ptttype.FILE_HEADER_RAW_SZ,
			},
			wantErr: false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			wg.Done()
			if err := DeleteRecord(tt.args.filename, tt.args.index, tt.args.theSize); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}
