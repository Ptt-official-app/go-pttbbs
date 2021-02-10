package api

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestCreateBoard(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testUserecRaw3)

	params0 := &CreateBoardParams{
		Brdname: "mnewtest",
	}

	path0 := &CreateBoardPath{
		ClsBid: 2,
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
		path       interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult interface{}
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				uuserID: "test",
				params:  params0,
				path:    path0,
			},
			expectedResult: testBoardSummary13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := CreateBoard(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotResult, tt.expectedResult)
		})
	}
}
