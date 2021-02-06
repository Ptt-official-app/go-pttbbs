package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const LOAD_AUTO_COMPLETE_BOARDS_R = "/boards/autocomplete"

type LoadAutoCompleteBoardsParams struct {
	StartIdx string `json:"start_idx,omitempty" form:"start_idx,omitempty" url:"start_idx,omitempty"`
	NBoards  int    `json:"max,omitempty" form:"max,omitempty" url:"max,omitempty"`
	Keyword  string `json:"brdname,omitempty" form:"brdname,omitempty" url:"brdname,omitempty"`
	Asc      bool   `json:"asc,omitempty" form:"asc,omitempty" url:"asc,omitempty"`
}

func LoadAutoCompleteBoardsWrapper(c *gin.Context) {
	params := &LoadAutoCompleteBoardsParams{}
	LoginRequiredQuery(LoadAutoCompleteBoards, params, c)
}

func LoadAutoCompleteBoards(remoteAddr string, uuserID bbs.UUserID, params interface{}) (interface{}, error) {
	theParams, ok := params.(*LoadAutoCompleteBoardsParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	summaries, nextIdx, err := bbs.LoadAutoCompleteBoards(uuserID, theParams.StartIdx, theParams.NBoards, theParams.Keyword, theParams.Asc)
	if err != nil {
		return nil, err
	}

	results := &LoadGeneralBoardsResult{
		Boards:  summaries,
		NextIdx: nextIdx,
	}

	return results, nil
}
