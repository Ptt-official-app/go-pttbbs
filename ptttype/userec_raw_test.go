package ptttype

import (
	"os"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestNewUserecRawWithFile(t *testing.T) {
	f, _ := os.Open("testcase/PASSWDS")
	defer f.Close()

	type args struct {
		file *os.File
	}
	tests := []struct {
		name     string
		args     args
		expected *UserecRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{f},
			expected: testUserecRaw1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserecRawWithFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserecRawWithFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			types.TDeepEqual(t, got, tt.expected)
		})
	}
}
