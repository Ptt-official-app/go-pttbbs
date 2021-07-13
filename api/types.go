package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

type ApiFunc func(remoteAddr string, params interface{}) (interface{}, error)

type LoginRequiredApiFunc func(remoteAddr string, uuserID bbs.UUserID, params interface{}) (interface{}, error)

type LoginRequiredPathApiFunc func(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (interface{}, error)

type JwtClaim struct {
	ClientInfo string `json:"cli"`
	UUserID    string `json:"sub"`
	Expire     int    `json:"exp"`
}

type EmailJwtClaim struct {
	ClientInfo string `json:"cli"`
	UUserID    string `json:"sub"`
	Email      string `json:"eml"`
	Expire     int    `json:"exp"`
	Context    string `json:"ctx"`
}

type errResult struct {
	Msg string
}
