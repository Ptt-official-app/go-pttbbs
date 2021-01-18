package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const LOAD_GENERAL_BOARDS_R = "/boards"

type LoadGeneralBoardsParams struct {
	StartIdx string `json:"start_idx,omitempty" form:"start_idx,omitempty" url:"start_idx,omitempty"`
	NBoards  int    `json:"max,omitempty" form:"max,omitempty" url:"max,omitempty"`
	Keyword  []byte `json:"keyword,omitempty" form:"keyword,omitempty" url:"keyword,omitempty"` //sending utf8-bytes from middleware
}

type LoadGeneralBoardsResult struct {
	Boards  []*bbs.BoardSummary `json:"data"`
	NextIdx string              `json:"next_idx"`
}

func LoadGeneralBoardsWrapper(c *gin.Context) {
	params := &LoadGeneralBoardsParams{}
	LoginRequiredQuery(LoadGeneralBoards, params, c)
}

func LoadGeneralBoards(remoteAddr string, uuserID bbs.UUserID, params interface{}) (interface{}, error) {
	return loadGeneralBoardsCore(remoteAddr, uuserID, params, ptttype.BSORT_BY_NAME)
}

func loadGeneralBoardsCore(remoteAddr string, uuserID bbs.UUserID, params interface{}, bsortBy ptttype.BSortBy) (interface{}, error) {
	theParams, ok := params.(*LoadGeneralBoardsParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	summary, nextIdx, err := bbs.LoadGeneralBoards(uuserID, theParams.StartIdx, theParams.NBoards, theParams.Keyword, bsortBy)
	if err != nil {
		return nil, err
	}

	results := &LoadGeneralBoardsResult{
		Boards:  summary,
		NextIdx: nextIdx,
	}

	return results, nil
}
