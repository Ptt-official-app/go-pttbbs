package api

import (
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

func LoginRequiredJSON(theFunc LoginRequiredAPIFunc, params interface{}, c *gin.Context) {
	err := c.ShouldBindJSON(params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	loginRequiredProcess(theFunc, params, c)
}

func LoginRequiredQuery(theFunc LoginRequiredAPIFunc, params interface{}, c *gin.Context) {
	err := c.ShouldBindQuery(params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	loginRequiredProcess(theFunc, params, c)
}

func loginRequiredProcess(theFunc LoginRequiredAPIFunc, params interface{}, c *gin.Context) {
	host := strings.TrimSpace(c.GetHeader("Host"))
	if !isValidHost(host) {
		processResult(c, nil, ErrInvalidHost)
		return
	}

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
	remoteAddr := strings.TrimSpace(c.GetHeader("X-Forwarded-For"))
	if !isValidRemoteAddr(remoteAddr) {
		processResult(c, nil, ErrInvalidRemoteAddr)
		return
	}

	jwt := GetJwt(c)

	userID, _, err := VerifyJwt(jwt)
	if err != nil {
		userID = bbs.UUserID(GUEST)
	}

	result, err := theFunc(remoteAddr, userID, params)
	processResult(c, result, err)
}
