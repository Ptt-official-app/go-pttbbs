package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const IS_BOARD_VALID_USER_R = "/board/:bid/isvalid"

type IsBoardValidUserPath struct {
	BoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

type IsBoardValidUserResult struct {
	IsValid bool `json:"is_valid"`
}

func IsBoardValidUserWrapper(c *gin.Context) {
	path := &IsBoardValidUserPath{}
	LoginRequiredPathQuery(IsBoardValidUser, nil, path, c)
}

func IsBoardValidUser(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	thePath, ok := path.(*IsBoardValidUserPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	isValid, err := bbs.IsBoardValidUser(uuserID, thePath.BoardID)
	if err != nil {
		return nil, err
	}

	return &IsBoardValidUserResult{IsValid: isValid}, nil
}
