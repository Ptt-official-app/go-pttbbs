package boardd

import (
	"context"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

//LoadHotBoards
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1125
func LoadHotBoards(user *ptttype.UserecRaw, uid ptttype.UID) (summary []*ptttype.BoardSummaryRaw, err error) {
	ctx := context.Background()
	req := &HotboardRequest{}
	resp, err := cli.Hotboard(ctx, req)
	if err != nil {
		return nil, err
	}

	summary = make([]*ptttype.BoardSummaryRaw, 0, len(resp.Boards))
	for _, each := range resp.Boards {
		eachSummary := boardToBoardSummaryRaw(each)
		summary = append(summary, eachSummary)
	}
	return summary, nil
}
