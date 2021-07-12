package ptttype

import (
	"os"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
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
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := NewUserecRawWithFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserecRawWithFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "userec", got, tt.expected)
		})
	}
	wg.Wait()
}
