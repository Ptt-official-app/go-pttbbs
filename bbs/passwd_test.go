package bbs

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestOpenUserecFile(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		filename string
	}
	tests := []struct {
		name     string
		args     args
		expected []*Userec
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{"testcase/passwd/01.PASSWDS"},
			expected: testOpenUserecFile1,
		},
		{
			args:     args{"testcase/passwd/01.PASSWDS.corrupt"},
			expected: testOpenUserecFile1[:TEST_N_OPEN_USER_FILE_1-1],
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenUserecFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenUserecFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("OpenUserecFile: got: %v expected: %v", got, tt.expected)
			}

			if len(got) != len(tt.expected) {
				t.Errorf("OpenUserecFile: len(got): %v expected: %v", len(got), len(tt.expected))
			}
			for idx, each := range tt.expected {
				if idx >= len(got) {
					break
				}
				types.TDeepEqual(t, got[idx], each)
			}
		})
	}
}
