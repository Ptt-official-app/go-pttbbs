package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const ATTEMPT_CHANGE_EMAIL_R = "/user/:uid/attemptchangeemail"

type AttemptChangeEmailParams struct {
	ClientInfo string `json:"client_info"`
	Passwd     string `json:"password"`
	Email      string `json:"email"`
}

type AttemptChangeEmailPath struct {
	UserID bbs.UUserID `uri:"uid" binding:"required"`
}

type AttemptChangeEmailResult struct {
	UserID bbs.UUserID `json:"user_id"`
	Jwt    string      `json:"email_token"`
}

func AttemptChangeEmailWrapper(c *gin.Context) {
	params := &AttemptChangeEmailParams{}
	path := &AttemptChangeEmailPath{}
	LoginRequiredPathJSON(AttemptChangeEmail, params, path, c)
}

func AttemptChangeEmail(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*AttemptChangeEmailParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*AttemptChangeEmailPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if !userInfoIsValidUser(uuserID, thePath.UserID) {
		return nil, ErrInvalidUser
	}

	err = bbs.CheckPasswd(uuserID, theParams.Passwd, remoteAddr)
	if err != nil {
		return nil, err
	}

	token, err := CreateEmailToken(thePath.UserID, theParams.ClientInfo, theParams.Email, CONTEXT_CHANGE_EMAIL)
	if err != nil {
		return nil, err
	}

	result = &AttemptChangeEmailResult{
		UserID: thePath.UserID,
		Jwt:    token,
	}

	return result, nil
}
