package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const LOAD_BOARD_DETAIL_R = "/board/:bid/detail"

type LoadBoardDetailParams struct{}

type LoadBoardDetailPath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

// XXX LoadBoardDetailResult is a pointer
//    It's possible that this is unavoidable,
//    and we need to change all the result-type to be pointer-based.
type LoadBoardDetailResult *bbs.BoardDetail

func LoadBoardDetailWrapper(c *gin.Context) {
	path := &LoadBoardDetailPath{}
	loginRequiredPathProcess(LoadBoardDetail, nil, path, c)
}

func LoadBoardDetail(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (ret interface{}, err error) {
	thePath, ok := path.(*LoadBoardDetailPath)
	if !ok {
		return nil, ErrInvalidPath
	}
	detail, err := bbs.LoadBoardDetail(uuserID, thePath.BBoardID)
	if err != nil {
		return nil, err
	}

	return LoadBoardDetailResult(detail), nil
}
