package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const CHANGE_PASSWD_R = "/users/:uid/changepasswd"

type ChangePasswdParams struct {
	ClientInfo string `json:"client_info"`
	OrigPasswd string `json:"orig_password"`
	Passwd     string `json:"password"`
}

type ChangePasswdPath struct {
	UserID bbs.UUserID `uri:"uid" binding:"required"`
}

type ChangePasswdResult struct {
	UserID    bbs.UUserID `json:"user_id"`
	Jwt       string      `json:"access_token"`
	TokenType string      `json:"token_type"`
}

func ChangePasswdWrapper(c *gin.Context) {
	params := &ChangePasswdParams{}
	path := &ChangePasswdPath{}
	LoginRequiredPathJSON(ChangePasswd, params, path, c)
}

func ChangePasswd(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*ChangePasswdParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*ChangePasswdPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if !userIsValidUser(uuserID, thePath.UserID) {
		return nil, ErrInvalidUser
	}

	err = bbs.ChangePasswd(uuserID, theParams.OrigPasswd, theParams.Passwd, remoteAddr)
	if err != nil {
		return nil, err
	}

	token, err := createToken(uuserID, theParams.ClientInfo)
	if err != nil {
		return nil, err
	}

	result = &ChangePasswdResult{
		UserID:    uuserID,
		Jwt:       token,
		TokenType: "bearer",
	}

	return result, nil
}
