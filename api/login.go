package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

const LOGIN_R = "/token"

type LoginParams struct {
	ClientInfo string `json:"client_info"`
	Username   string `json:"username"`
	Passwd     string `json:"password"`
}

type LoginResult struct {
	UserID        bbs.UUserID `json:"user_id"`
	Jwt           string      `json:"access_token"`
	TokenType     string      `json:"token_type"`
	Refresh       string      `json:"refresh_token"`
	AccessExpire  types.Time4 `json:"access_expire"`
	RefreshExpire types.Time4 `json:"refresh_expire"`
}

func LoginWrapper(c *gin.Context) {
	params := &LoginParams{}
	JSON(Login, params, c)
}

func Login(remoteAddr string, params interface{}, c *gin.Context) (interface{}, error) {
	loginParams, ok := params.(*LoginParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	uuserID, err := bbs.Login(loginParams.Username, loginParams.Passwd, remoteAddr)
	if err != nil {
		return nil, ErrLoginFailed
	}

	token, accessExpireTime, err := CreateToken(uuserID, loginParams.ClientInfo)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpireTime, err := CreateRefreshToken(uuserID, loginParams.ClientInfo)
	if err != nil {
		return nil, err
	}

	result := &LoginResult{
		UserID:        uuserID,
		Jwt:           token,
		TokenType:     "bearer",
		Refresh:       refreshToken,
		AccessExpire:  accessExpireTime,
		RefreshExpire: refreshExpireTime,
	}

	return result, nil
}
