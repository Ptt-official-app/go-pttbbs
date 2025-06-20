package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

// XXX LoadBoardSummaryResult is a pointer
//
//	It's possible that this is unavoidable,
//	and we need to change all the result-type to be pointer-based.
func LoadBoardSummaryAllGuestWrapper(c *gin.Context) {
	path := &LoadBoardSummaryPath{}
	PathQuery(LoadBoardSummaryAllGuest, nil, path, c)
}

func LoadBoardSummaryAllGuest(remoteAddr string, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	thePath, ok := path.(*LoadBoardSummaryPath)
	if !ok {
		return nil, ErrInvalidPath
	}
	summary, err := bbs.LoadBoardSummaryAllGuest(thePath.BBoardID)
	if summary == nil {
		return nil, ErrInvalidParams
	}
	if err != nil {
		return nil, err
	}
	result = LoadBoardSummaryResult(summary)

	return result, nil
}
