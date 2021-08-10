package ptt

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func Test_isBannedByBoard(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	nowTS := types.NowTS()
	banTS := nowTS + 3600

	type args struct {
		user    *ptttype.UserecRaw
		board   *ptttype.BoardHeaderRaw
		isToBan bool
	}
	tests := []struct {
		name             string
		args             args
		expectedExpireTS types.Time4
		expectedReason   string
	}{
		// TODO: Add test cases.
		{
			args:             args{user: testUserecRaw1, board: testBoardHeaderRaw2},
			expectedExpireTS: 0,
			expectedReason:   "",
		},
		{
			args:             args{user: testUserecRaw1, board: testBoardHeaderRaw2, isToBan: true},
			expectedExpireTS: banTS,
			expectedReason:   "test",
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if tt.args.isToBan {
				filename, _ := bakumanMakeTagFilename(&tt.args.user.UserID, types.Cstr(tt.args.board.Brdname[:]), BAKUMAN_OBJECT_TYPE_BOARD, true)
				file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o600)
				if err != nil {
					t.Errorf("unable to open file: filename: %v e: %v", filename, err)
					return
				}
				defer file.Close()

				write := bufio.NewWriter(file)
				fmt.Fprintf(write, "%v\ntest", banTS)
				write.Flush()
				logrus.Infof("isToBan: filename: %v", filename)
			}

			gotExpireTS, gotReason := isBannedByBoard(tt.args.user, tt.args.board)

			if !reflect.DeepEqual(gotExpireTS, tt.expectedExpireTS) {
				t.Errorf("isBannedByBoard() gotExpireTS = %v, want %v", gotExpireTS, tt.expectedExpireTS)
			}
			if gotReason != tt.expectedReason {
				t.Errorf("isBannedByBoard() gotReason = %v, want %v", gotReason, tt.expectedReason)
			}
		})
		wg.Wait()
	}
}
