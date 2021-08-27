package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const IS_BOARDS_VALID_USER_R = "/boards/isvalid"

type IsBoardsValidUserParams struct {
	BoardIDs []bbs.BBoardID `uri:"bids" binding:"required"`
}

type IsBoardsValidUserResult struct {
	IsValid map[bbs.BBoardID]bool `json:"is_valid"`
}

func IsBoardsValidUserWrapper(c *gin.Context) {
	params := &IsBoardsValidUserParams{}
	LoginRequiredJSON(IsBoardsValidUser, params, c)
}

func IsBoardsValidUser(remoteAddr string, uuserID bbs.UUserID, params interface{}) (result interface{}, err error) {
	theParams, ok := params.(*IsBoardsValidUserParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	isValid, err := bbs.IsBoardsValidUser(uuserID, theParams.BoardIDs)
	if err != nil {
		return nil, err
	}

	return &IsBoardsValidUserResult{IsValid: isValid}, nil
}
