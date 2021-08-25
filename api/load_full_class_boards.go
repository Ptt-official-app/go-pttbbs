package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const LOAD_FULL_CLASS_BOARDS_R = "/cls/boards"

type LoadFullClassBoardsParams struct {
	StartBid ptttype.Bid `json:"start_bid,omitempty" form:"start_bid,omitempty" url:"start_bid,omitempty"`
	NBoards  int         `json:"max" form:"max" url:"max"`
	IsSystem bool        `json:"system,omitempty" form:"system,omitempty" url:"system"`
}

type LoadFullClassBoardsResult struct {
	Boards  []*bbs.BoardSummary `json:"data"`
	NextBid ptttype.Bid         `json:"next_bid"`
}

func LoadFullClassBoardsWrapper(c *gin.Context) {
	params := &LoadFullClassBoardsParams{}
	LoginRequiredQuery(LoadFullClassBoards, params, c)
}

func LoadFullClassBoards(remoteAddr string, uuserID bbs.UUserID, params interface{}) (ret interface{}, err error) {
	theParams, ok := params.(*LoadFullClassBoardsParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	if !theParams.IsSystem {
		return nil, ErrInvalidUser
	}

	uuserID = bbs.UUserID(string(ptttype.STR_SYSOP))

	summaries, nextBid, err := bbs.LoadFullClassBoards(uuserID, theParams.StartBid, theParams.NBoards)
	if err != nil {
		return nil, err
	}

	ret = &LoadFullClassBoardsResult{
		Boards:  summaries,
		NextBid: nextBid,
	}

	return ret, nil
}
