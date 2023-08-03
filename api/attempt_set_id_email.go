package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const ATTEMPT_SET_ID_EMAIL_R = "/user/:uid/attemptsetidemail"

type AttemptSetIDEmailParams struct {
	ClientInfo string `json:"client_info"`
	Passwd     string `json:"password"`
	Email      string `json:"email"`
}

type AttemptSetIDEmailPath struct {
	UserID bbs.UUserID `uri:"uid" binding:"required"`
}

type AttemptSetIDEmailResult struct {
	UserID bbs.UUserID `json:"user_id"`
	Jwt    string      `json:"email_token"`
}

func AttemptSetIDEmailWrapper(c *gin.Context) {
	params := &AttemptSetIDEmailParams{}
	path := &AttemptSetIDEmailPath{}
	LoginRequiredPathJSON(AttemptSetIDEmail, params, path, c)
}

func AttemptSetIDEmail(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*AttemptSetIDEmailParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*AttemptSetIDEmailPath)
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

	err = bbs.IsValidIDEmail(theParams.Email)
	if err != nil {
		return nil, ErrInvalidIDEmail
	}

	token, err := CreateEmailToken(thePath.UserID, theParams.ClientInfo, theParams.Email, CONTEXT_SET_ID_EMAIL)
	if err != nil {
		return nil, err
	}

	result = &AttemptSetIDEmailResult{
		UserID: thePath.UserID,
		Jwt:    token,
	}

	return result, nil
}
