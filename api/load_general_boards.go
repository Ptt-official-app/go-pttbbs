package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

const LOAD_GENERAL_BOARDS_R = "/boards"

type LoadGeneralBoardsParams struct {
	StartIdx string `json:"start_idx,omitempty" form:"start_idx,omitempty" url:"start_idx,omitempty"`
	NBoards  int    `json:"max,omitempty" form:"max,omitempty" url:"max,omitempty"`
	Keyword  string `json:"keyword,omitempty" form:"keyword,omitempty" url:"keyword,omitempty"`
}

type LoadGeneralBoardsResult struct {
	Boards  []*bbs.BoardSummary `json:"data"`
	NextIdx string              `json:"next_idx"`
}

func LoadGeneralBoards(remoteAddr string, uuserID bbs.UUserID, params interface{}) (interface{}, error) {
	loadGeneralBoardsParams, ok := params.(*LoadGeneralBoardsParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	summary, nextIdx, err := bbs.LoadGeneralBoards(uuserID, loadGeneralBoardsParams.StartIdx, loadGeneralBoardsParams.NBoards, loadGeneralBoardsParams.Keyword)
	if err != nil {
		return nil, err
	}

	results := &LoadGeneralBoardsResult{
		Boards:  summary,
		NextIdx: nextIdx,
	}

	return results, nil
}
