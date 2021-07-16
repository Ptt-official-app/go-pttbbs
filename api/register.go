package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const REGISTER_R = "/register"

type RegisterParams struct {
	ClientInfo string `json:"client_info"`
	Username   string `json:"username"`
	Passwd     string `json:"password"`
	Email      string `json:"email,omitempty"`

	Nickname []byte `json:"nickname,omitempty"` // sending utf8-bytes from middleware.
	Realname []byte `json:"realname,omitempty"`
	Career   []byte `json:"career,omitempty"`
	Address  []byte `json:"address,omitempty"`
	Over18   bool   `json:"over18"`
}

type RegisterResult struct {
	UserID    bbs.UUserID `json:"user_id"`
	Jwt       string      `json:"access_token"`
	TokenType string      `json:"token_type"`
}

func RegisterWrapper(c *gin.Context) {
	params := &RegisterParams{}
	JSON(Register, params, c)
}

func Register(remoteAddr string, params interface{}) (interface{}, error) {
	registerParams, ok := params.(*RegisterParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	uuserID, err := bbs.Register(
		registerParams.Username,
		registerParams.Passwd,
		remoteAddr,
		registerParams.Email,

		registerParams.Nickname,
		registerParams.Realname,
		registerParams.Career,
		registerParams.Address,
		registerParams.Over18,
	)
	if err != nil {
		return nil, err
	}

	token, err := CreateToken(uuserID, registerParams.ClientInfo)
	if err != nil {
		return nil, err
	}

	result := &RegisterResult{
		UserID:    uuserID,
		Jwt:       token,
		TokenType: "bearer",
	}

	return result, nil
}
