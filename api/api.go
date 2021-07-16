package api

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func JSON(theFunc APIFunc, params interface{}, c *gin.Context) {
	err := c.ShouldBindJSON(params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	process(theFunc, params, c)
}

func Query(theFunc APIFunc, params interface{}, c *gin.Context) {
	err := c.ShouldBindQuery(params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	process(theFunc, params, c)
}

func process(theFunc APIFunc, params interface{}, c *gin.Context) {
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

	result, err := theFunc(remoteAddr, params)
	processResult(c, result, err)
}
