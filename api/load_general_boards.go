package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

type LoadGeneralBoardsParams struct {
	StartIdx string `json:"start_idx"`
	NBoards  int    `json:"max"`
	Keyword  string `json:"keyword"`
}

type LoadGeneralBoardsResult struct {
	Boards  []*bbs.BoardSummary `json:"boards"`
	NextIdx string              `json:"next_bid"`
}

func LoadGeneralBoards(remoteAddr string, userID string, params interface{}) (interface{}, error) {
	loadGeneralBoardsParams, ok := params.(*LoadGeneralBoardsParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	summary, nextIdx, err := bbs.LoadGeneralBoards(userID, loadGeneralBoardsParams.StartIdx, loadGeneralBoardsParams.NBoards, loadGeneralBoardsParams.Keyword)
	if err != nil {
		return nil, err
	}

	results := &LoadGeneralBoardsResult{
		Boards:  summary,
		NextIdx: nextIdx,
	}

	return results, nil
}
