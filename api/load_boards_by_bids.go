package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const LOAD_BOARDS_BY_BIDS_R = "/boards/bybids"

type LoadBoardsByBidsParams struct {
	Bids []ptttype.Bid `json:"bids"`
}

type LoadBoardsByBidsResult struct {
	Boards []*bbs.BoardSummary `json:"data"`
}

func LoadBoardsByBidsWrapper(c *gin.Context) {
	params := &LoadBoardsByBidsParams{}
	LoginRequiredJSON(LoadBoardsByBids, params, c)
}

func LoadBoardsByBids(remoteAddr string, uuserID bbs.UUserID, params interface{}) (result interface{}, err error) {
	theParams, ok := params.(*LoadBoardsByBidsParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	summaries, err := bbs.LoadBoardsByBids(uuserID, theParams.Bids)
	if err != nil {
		return nil, err
	}

	result = &LoadBoardsByBidsResult{
		Boards: summaries,
	}

	return result, nil
}
