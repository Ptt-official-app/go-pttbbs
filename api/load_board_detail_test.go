package api

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadBoardDetail(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	params := &LoadBoardDetailParams{}

	path := &LoadBoardDetailPath{
		BBoardID: "6_ALLPOST",
	}

	expected := testBoardDetail6

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
		path       interface{}
	}
	tests := []struct {
		name            string
		args            args
		expectedResults LoadBoardDetailResult
		wantErr         bool
	}{
		// TODO: Add test cases & more field.
		{
			args:            args{remoteAddr: testIP, uuserID: "SYSOP", params: params, path: path},
			expectedResults: expected,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResults, err := LoadBoardDetail(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBoardDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result, _ := gotResults.(LoadBoardDetailResult)
			testutil.TDeepEqual(t, "got", result, tt.expectedResults)
		})
		wg.Wait()
	}
}
