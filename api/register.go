package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

const REGISTER_R = "/register"

type RegisterParams struct {
	Username string `json:"username"`
	Passwd   string `json:"password"`
	Email    string `json:"email,omitempty"`

	Nickname []byte `json:"nickname,omitempty"` //sending utf8-bytes from middleware.
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

func Register(remoteAddr string, params interface{}) (interface{}, error) {
	registerParams, ok := params.(*RegisterParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	user, err := bbs.Register(
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

	token, err := createToken(user)
	if err != nil {
		return nil, err
	}

	result := &RegisterResult{
		UserID:    user.UUserID,
		Jwt:       token,
		TokenType: "bearer",
	}

	return result, nil
}
