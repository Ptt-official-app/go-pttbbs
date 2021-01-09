package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const LOAD_BOARD_SUMMARY_R = "/boards/:bid/summary"

type LoadBoardSummaryParams struct{}

type LoadBoardSummaryPath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

//XXX LoadBoardSummaryResult is a pointer
//    It's possible that this is unavoidable,
//    and we need to change all the result-type to be pointer-based.
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
