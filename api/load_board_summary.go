package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

const LOAD_BOARD_SUMMARY_R = "/boards/:bid/summary"

type LoadBoardSummaryPath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

type LoadBoardSummaryResult *bbs.BoardSummary

func LoadBoardSummary(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (interface{}, error) {
	thePath, ok := path.(*LoadBoardSummaryPath)
	if !ok {

		return nil, ErrInvalidPath
	}
	summary, err := bbs.LoadBoardSummary(uuserID, thePath.BBoardID)
	if err != nil {
		return nil, err
	}
	results := LoadBoardSummaryResult(summary)

	return results, nil
}
