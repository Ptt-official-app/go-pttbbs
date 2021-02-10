package bbs

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestCreateBoard(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testUserecRaw3)

	type args struct {
		userID       UUserID
		clsBid       ptttype.Bid
		brdname      string
		brdClass     []byte
		brdTitle     []byte
		BMs          []UUserID
		brdAttr      ptttype.BrdAttr
		level        ptttype.PERM
		chessCountry ptttype.ChessCode
		isGroup      bool
	}
	tests := []struct {
		name            string
		args            args
		expectedBoardID BBoardID
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				userID:   "SYSOP",
				clsBid:   5,
				brdname:  "mnewboard",
				brdClass: []byte("CPBL"),
				brdTitle: []byte("new-board"),
			},
			wantErr: true,
		},
		{
			args: args{
				userID:   "test",
				clsBid:   5,
				brdname:  "mnewboard",
				brdClass: []byte("CPBL"),
				brdTitle: []byte("new-board"),
			},
			expectedBoardID: "13_mnewboard",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBoardID, err := CreateBoard(tt.args.userID, tt.args.clsBid, tt.args.brdname, tt.args.brdClass, tt.args.brdTitle, tt.args.BMs, tt.args.brdAttr, tt.args.level, tt.args.chessCountry, tt.args.isGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBoardID != tt.expectedBoardID {
				t.Errorf("CreateBoard() = %v, want %v", gotBoardID, tt.expectedBoardID)
			}
		})
	}
}
