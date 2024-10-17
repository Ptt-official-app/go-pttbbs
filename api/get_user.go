package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const GET_USER_R = "/user/:uid/information"

type GetUserParams struct{}

type GetUserPath struct {
	UserID bbs.UUserID `uri:"uid" binding:"required"`
}

type GetUserResult *bbs.Userec

func GetUserWrapper(c *gin.Context) {
	path := &GetUserPath{}

	LoginRequiredPathQuery(GetUser, nil, path, c)
}

func GetUser(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	thePath, ok := path.(*GetUserPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if !userInfoIsValidUser(uuserID, thePath.UserID) {
		return nil, ErrInvalidUser
	}

	user, err := bbs.GetUser(thePath.UserID)
	if err != nil {
		return nil, err
	}

	return GetUserResult(user), nil
}
