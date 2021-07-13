package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestIndex(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())
	// setupTest moves inside for-loop
	// teardownTest moves inside for-loop

	type args struct {
		uuserID bbs.UUserID
		params  interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{"1_SYSOP", nil},
			expected: &IndexResult{Data: "index"},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			got, err := Index(testIP, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Index() = %v, want %v", got, tt.expected)
			}
		})
		wg.Wait()
	}
}
