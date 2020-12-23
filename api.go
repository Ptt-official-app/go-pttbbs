package main

import (
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

type Api struct {
	Func   api.ApiFunc
	Params interface{}
}

type LoginRequiredApi struct {
	Func   api.LoginRequiredApiFunc
	Params interface{}
}

type LoginRequiredPathApi struct {
	Func   api.LoginRequiredPathApiFunc
	Path   interface{}
	Params interface{}
}

func NewApi(f api.ApiFunc, params interface{}) *Api {
	return &Api{Func: f, Params: params}
}

func (a *Api) Json(c *gin.Context) {
	err := c.ShouldBindJSON(a.Params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	a.process(c)
}

func (a *Api) Query(c *gin.Context) {
	err := c.ShouldBindQuery(a.Params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	a.process(c)
}

func (a *Api) process(c *gin.Context) {
	host := strings.TrimSpace(c.GetHeader("Host"))
	if !isValidHost(host) {
		processResult(c, nil, ErrInvalidHost)
		return
	}

	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
	remoteAddr := strings.TrimSpace(c.GetHeader("X-Forwarded-For"))
	if !isValidRemoteAddr(remoteAddr) {
		processResult(c, nil, ErrInvalidRemoteAddr)
		return
	}

	result, err := a.Func(remoteAddr, a.Params)
	processResult(c, result, err)
}

func NewLoginRequiredApi(f api.LoginRequiredApiFunc, params interface{}) *LoginRequiredApi {
	return &LoginRequiredApi{Func: f, Params: params}
}

func (a *LoginRequiredApi) Json(c *gin.Context) {
	err := c.ShouldBindJSON(a.Params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	a.process(c)
}

func (a *LoginRequiredApi) Query(c *gin.Context) {
	err := c.ShouldBindQuery(a.Params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	a.process(c)
}

func (a *LoginRequiredApi) process(c *gin.Context) {
	host := strings.TrimSpace(c.GetHeader("Host"))
	if !isValidHost(host) {
		processResult(c, nil, ErrInvalidHost)
		return
	}

	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
	remoteAddr := strings.TrimSpace(c.GetHeader("X-Forwarded-For"))
	if !isValidRemoteAddr(remoteAddr) {
		processResult(c, nil, ErrInvalidRemoteAddr)
		return
	}

	tokenStr := strings.TrimSpace(c.GetHeader("Authorization"))
	tokenList := strings.Split(tokenStr, " ")
	if len(tokenList) != 2 {
		processResult(c, nil, ErrInvalidToken)
		return
	}
	jwt := tokenList[1]

	userID, err := api.VerifyJwt(jwt)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	result, err := a.Func(remoteAddr, userID, a.Params)
	processResult(c, result, err)

}

func NewLoginRequiredPathApi(f api.LoginRequiredPathApiFunc, params interface{}, path interface{}) *LoginRequiredPathApi {
	return &LoginRequiredPathApi{Func: f, Params: params, Path: path}
}

func (a *LoginRequiredPathApi) Json(c *gin.Context) {
	err := c.ShouldBindJSON(a.Params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	a.process(c)
}

func (a *LoginRequiredPathApi) Query(c *gin.Context) {
	err := c.ShouldBindQuery(a.Params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	a.process(c)
}

func (a *LoginRequiredPathApi) process(c *gin.Context) {
	err := c.ShouldBindUri(a.Path)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	host := strings.TrimSpace(c.GetHeader("Host"))
	if !isValidHost(host) {
		processResult(c, nil, ErrInvalidHost)
		return
	}

	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
	remoteAddr := strings.TrimSpace(c.GetHeader("X-Forwarded-For"))
	if !isValidRemoteAddr(remoteAddr) {
		processResult(c, nil, ErrInvalidRemoteAddr)
		return
	}

	tokenStr := strings.TrimSpace(c.GetHeader("Authorization"))
	tokenList := strings.Split(tokenStr, " ")
	if len(tokenList) != 2 {
		processResult(c, nil, ErrInvalidToken)
		return
	}
	jwt := tokenList[1]

	userID, err := api.VerifyJwt(jwt)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	result, err := a.Func(remoteAddr, userID, a.Params, a.Path)
	processResult(c, result, err)

}

func processResult(c *gin.Context, result interface{}, err error) {
	switch err {
	case nil:
		c.JSON(200, result)

	case ErrInvalidHost:
		c.JSON(400, &errResult{err.Error()})
	case ErrInvalidRemoteAddr:
		c.JSON(400, &errResult{err.Error()})

	case api.ErrInvalidParams:
		c.JSON(400, &errResult{err.Error()})
	case api.ErrInvalidPath:
		c.JSON(400, &errResult{err.Error()})

	case bbs.ErrInvalidParams:
		c.JSON(400, &errResult{err.Error()})

	case ptttype.ErrUserIDAlreadyExists:
		c.JSON(400, &errResult{err.Error()})

	//401
	case ErrInvalidToken:
		c.JSON(401, &errResult{err.Error()})

	case api.ErrInvalidToken:
		c.JSON(401, &errResult{err.Error()})
	case api.ErrLoginFailed:
		c.JSON(401, &errResult{err.Error()})

	default:
		c.JSON(500, &errResult{err.Error()})
	}
}
