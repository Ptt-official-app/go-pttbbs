package ptt

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func Test_boardPermStatNormally(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	type args struct {
		user  *ptttype.UserecRaw
		uid   ptttype.Uid
		board *ptttype.BoardHeaderRaw
		bid   ptttype.Bid
	}
	tests := []struct {
		name     string
		args     args
		expected ptttype.BoardStatAttr
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:  testUserecRaw1,
				uid:   1,
				board: testBoardHeaderRaw1,
				bid:   5,
			},
			expected: ptttype.NBRD_BOARD,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := boardPermStatNormally(tt.args.user, tt.args.uid, tt.args.board, tt.args.bid); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("boardPermStatNormally() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsBoardValidUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	type args struct {
		user    *ptttype.UserecRaw
		uid     ptttype.Uid
		boardID *ptttype.BoardID_t
		bid     ptttype.Bid
	}
	tests := []struct {
		name            string
		args            args
		expectedIsValid bool
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:    testUserecRaw1,
				uid:     1,
				boardID: &testBoardHeaderRaw1.Brdname,
				bid:     5,
			},
			expectedIsValid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsValid, err := IsBoardValidUser(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsBoardValidUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsValid != tt.expectedIsValid {
				t.Errorf("IsBoardValidUser() = %v, want %v", gotIsValid, tt.expectedIsValid)
			}
		})
	}
}
