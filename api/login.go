package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const LOGIN_R = "/token"

type LoginParams struct {
	ClientInfo string `json:"client_info"`
	Username   string `json:"username"`
	Passwd     string `json:"password"`
}

type LoginResult struct {
	UserID    bbs.UUserID `json:"user_id"`
	Jwt       string      `json:"access_token"`
	TokenType string      `json:"token_type"`
}

func LoginWrapper(c *gin.Context) {
	params := &LoginParams{}
	JSON(Login, params, c)
}

func Login(remoteAddr string, params interface{}) (interface{}, error) {
	loginParams, ok := params.(*LoginParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	user, err := bbs.Login(loginParams.Username, loginParams.Passwd, remoteAddr)
	if err != nil {
		return nil, ErrLoginFailed
	}

	token, err := createToken(user.UUserID, loginParams.ClientInfo)
	if err != nil {
		return nil, err
	}

	result := &LoginResult{
		UserID:    user.UUserID,
		Jwt:       token,
		TokenType: "bearer",
	}

	return result, nil
}
