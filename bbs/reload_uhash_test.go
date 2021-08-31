package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestReloadUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := ReloadUHash(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("ReloadUHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}
