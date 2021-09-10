package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadBoardDetail(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID  UUserID
		bboardID BBoardID
	}
	tests := []struct {
		name           string
		args           args
		expectedDetail *BoardDetail
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{uuserID: "SYSOP", bboardID: "1_SYSOP"},
			expectedDetail: testBoardDetail3,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotDetail, err := LoadBoardDetail(tt.args.uuserID, tt.args.bboardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotDetail, tt.expectedDetail)
		})
	}
	wg.Wait()
}
