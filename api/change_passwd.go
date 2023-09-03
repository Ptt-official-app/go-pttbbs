package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const CHANGE_PASSWD_R = "/user/:uid/changepasswd" //nolint

type ChangePasswdParams struct {
	ClientInfo string `json:"client_info"`
	OrigPasswd string `json:"orig_password"`
	Passwd     string `json:"password"`
}

type ChangePasswdPath struct {
	UserID bbs.UUserID `uri:"uid" binding:"required"`
}

type ChangePasswdResult LoginResult

func ChangePasswdWrapper(c *gin.Context) {
	params := &ChangePasswdParams{}
	path := &ChangePasswdPath{}
	LoginRequiredPathJSON(ChangePasswd, params, path, c)
}

func ChangePasswd(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*ChangePasswdParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*ChangePasswdPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if !userInfoIsValidUser(uuserID, thePath.UserID) {
		return nil, ErrInvalidUser
	}

	err = bbs.ChangePasswd(uuserID, theParams.OrigPasswd, theParams.Passwd, remoteAddr)
	if err != nil {
		return nil, err
	}

	token, accessExpireTime, err := CreateToken(uuserID, theParams.ClientInfo)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpireTime, err := CreateRefreshToken(uuserID, theParams.ClientInfo)
	if err != nil {
		return nil, err
	}

	result = &ChangePasswdResult{
		UserID:        uuserID,
		Jwt:           token,
		TokenType:     "bearer",
		Refresh:       refreshToken,
		AccessExpire:  accessExpireTime,
		RefreshExpire: refreshExpireTime,
	}

	return result, nil
}
