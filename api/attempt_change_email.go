package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const ATTEMPT_CHANGE_EMAIL_R = "/users/:uid/attemptchangeemail"

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

func AttemptChangeEmail(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*AttemptChangeEmailParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*AttemptChangeEmailPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if uuserID != thePath.UserID {
		return nil, ErrInvalidUser
	}

	err = bbs.CheckPasswd(uuserID, theParams.Passwd, remoteAddr)
	if err != nil {
		return nil, err
	}

	token, err := createEmailToken(thePath.UserID, theParams.ClientInfo, theParams.Email)
	if err != nil {
		return nil, err
	}

	result = &AttemptChangeEmailResult{
		UserID: uuserID,
		Jwt:    token,
	}

	return result, nil
}
