package api

import (
	"time"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const LOGIN_R = "/token"

type LoginParams struct {
	UserID string `json:"username"`
	Passwd string `json:"password"`
}

type LoginResult struct {
	Jwt       string `json:"access_token"`
	TokenType string `json:"token_type"`
}

func Login(remoteAddr string, params interface{}) (interface{}, error) {
	loginParams, ok := params.(*LoginParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	user, err := bbs.Login(loginParams.UserID, loginParams.Passwd, remoteAddr)
	if err != nil {
		return nil, ErrLoginFailed
	}

	token, err := createToken(user)
	if err != nil {
		return nil, err
	}

	result := &LoginResult{
		Jwt:       token,
		TokenType: "bearer",
	}

	return result, nil
}

func createToken(userec *bbs.Userec) (string, error) {
	var err error

	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: JWT_SECRET}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}

	cl := &JwtClaim{
		UserID: userec.Userid,
		Expire: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}

	raw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		return "", err
	}

	return raw, nil
}
