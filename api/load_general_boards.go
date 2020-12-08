package api

import "github.com/Ptt-official-app/go-pttbbs/bbs"

type LoadGeneralBoardsParams struct {
	StartIdx int32
	NBoards  int
	Keyword  []byte
}

type LoadGeneralBoardsResult struct {
	Boards  []*bbs.BoardSummary
	NextIdx int32
}

func LoadGeneralBoards(userID string, params interface{}) (interface{}, error) {
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
