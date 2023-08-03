package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

type APIFunc func(remoteAddr string, params interface{}, c *gin.Context) (interface{}, error)

type LoginRequiredAPIFunc func(remoteAddr string, uuserID bbs.UUserID, params interface{}, c *gin.Context) (interface{}, error)

type LoginRequiredPathAPIFunc func(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (interface{}, error)

type JwtClaim struct {
	ClientInfo string `json:"cli"`
	UUserID    string `json:"sub"`
	Expire     int    `json:"exp"`
}

type RefreshJwtClaim struct {
	ClientInfo string `json:"cli"`
	UUserID    string `json:"sub"`
	Expire     int    `json:"exp"`
	TheType    string `json:"typ"`
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
