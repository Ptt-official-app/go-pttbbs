package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const CHECK_EXISTS_USER_R = "/existsuser"

type CheckExistsUserParams struct {
	Username string `json:"username"`
}

type CheckExistsUserResult struct {
	UserID   bbs.UUserID `json:"user_id"`
	IsExists bool        `json:"is_exists"`
}

func CheckExistsUserWrapper(c *gin.Context) {
	params := &CheckExistsUserParams{}
	JSON(CheckExistsUser, params, c)
}

func CheckExistsUser(remoteAddr string, params interface{}) (result interface{}, err error) {
	theParams, ok := params.(*CheckExistsUserParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	uuserID, err := bbs.CheckExistsUser(theParams.Username)
	if err != nil {
		return nil, err
	}

	isExists := uuserID != ""

	return &CheckExistsUserResult{
		UserID:   uuserID,
		IsExists: isExists,
	}, nil
}
