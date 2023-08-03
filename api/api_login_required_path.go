package api

import (
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

func LoginRequiredPathJSON(theFunc LoginRequiredPathAPIFunc, params interface{}, path interface{}, c *gin.Context) {
	err := c.ShouldBindJSON(params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	loginRequiredPathProcess(theFunc, params, path, c)
}

func LoginRequiredPathQuery(theFunc LoginRequiredPathAPIFunc, params interface{}, path interface{}, c *gin.Context) {
	err := c.ShouldBindQuery(params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	loginRequiredPathProcess(theFunc, params, path, c)
}

func loginRequiredPathProcess(theFunc LoginRequiredPathAPIFunc, params interface{}, path interface{}, c *gin.Context) {
	err := c.ShouldBindUri(path)
	if err != nil {
		processResult(c, nil, err)
		return
	}

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

	userID, _, _, err := VerifyJwt(jwt, true)
	if err != nil {
		userID = bbs.UUserID(GUEST)
	}

	result, err := theFunc(remoteAddr, userID, params, path, c)
	processResult(c, result, err)
}
