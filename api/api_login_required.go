package api

import (
	"strings"

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

	userID, _ := apiGetUser(c)

	result, err := theFunc(remoteAddr, userID, params, c)
	processResult(c, result, err)
}
