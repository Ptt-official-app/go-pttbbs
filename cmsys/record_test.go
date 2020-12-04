package cmsys

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
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
