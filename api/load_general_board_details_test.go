package api

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestLoadGeneralBoardDetails(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	params0 := &LoadGeneralBoardDetailsParams{
		StartIdx: "",
		NBoards:  4,
		Asc:      true,
	}

	expected0 := &LoadGeneralBoardDetailsResult{
		Boards:  []*bbs.BoardDetail{testClassDetail2, testClassDetail5, testBoardDetail12, testBoardDetail6},
		NextIdx: "deleted",
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{uuserID: "SYSOP", params: params0},
			want: expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadGeneralBoardDetails(tt.args.remoteAddr, tt.args.uuserID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralBoardDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "got", got, tt.want)
		})
	}
}
