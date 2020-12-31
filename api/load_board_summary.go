package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const LOAD_BOARD_SUMMARY_R = "/boards/:bid/summary"

type LoadBoardSummaryPath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

type LoadBoardSummaryResult *bbs.BoardSummary

func LoadBoardSummaryWrapper(c *gin.Context) {

	path := &LoadBoardSummaryPath{}
	loginRequiredPathProcess(LoadBoardSummary, nil, path, c)

}

func LoadBoardSummary(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (results interface{}, err error) {
	thePath, ok := path.(*LoadBoardSummaryPath)
	if !ok {

		return nil, ErrInvalidPath
	}
	summary, err := bbs.LoadBoardSummary(uuserID, thePath.BBoardID)
	if summary == nil {
		return nil, ErrInvalidParams
	}
	if err != nil {
		return nil, err
	}
	results = LoadBoardSummaryResult(summary)

	return results, nil
}
