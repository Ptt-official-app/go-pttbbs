package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"gopkg.in/square/go-jose.v2/jwt"
)

type ApiFunc func(remoteAddr string, params interface{}) (interface{}, error)

type LoginRequiredApiFunc func(remoteAddr string, uuserID bbs.UUserID, params interface{}) (interface{}, error)

type LoginRequiredPathApiFunc func(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (interface{}, error)

type JwtClaim struct {
	UUserID bbs.UUserID
	Expire  *jwt.NumericDate
}

type errResult struct {
	Msg string
}
