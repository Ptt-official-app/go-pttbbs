package main

import (
	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Api struct {
	Func   api.ApiFunc
	Params interface{}
}

type LoginRequiredParams struct {
	UserID string `json:"u"`
	Jwt    string `json:"j"`
	Data   interface{}
}

type LoginRequiredApi struct {
	Func   api.LoginRequiredApiFunc
	Params interface{}
}

func NewApi(f api.ApiFunc, params interface{}) *Api {
	return &Api{Func: f, Params: params}
}

func NewLoginRequiredApi(f api.LoginRequiredApiFunc, params interface{}) *LoginRequiredApi {
	return &LoginRequiredApi{Func: f, Params: params}
}

func (api *Api) Json(c *gin.Context) {
	err := c.ShouldBindJSON(api.Params)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	result, err := api.Func(api.Params)
	processResult(c, result, err)
}

func (api *LoginRequiredApi) LoginRequiredJson(c *gin.Context) {
	loginParams := &LoginRequiredParams{Data: api.Params}
	err := c.ShouldBindJSON(loginParams)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	err = verifyJwt(loginParams.UserID, loginParams.Jwt)
	if err != nil {
		processResult(c, nil, err)
		return
	}

	result, err := api.Func(loginParams.UserID, loginParams.Data)
	processResult(c, result, err)
}

func verifyJwt(userID string, raw string) error {
	tok, err := jwt.ParseSigned(raw)
	if err != nil {
		return ErrInvalidToken
	}

	cl := &api.JwtClaim{}
	if err := tok.Claims(api.JWT_SECRET, cl); err != nil {
		return ErrInvalidToken
	}

	if cl.UserID != userID {
		return ErrInvalidToken
	}

	return nil
}

func processResult(c *gin.Context, result interface{}, err error) {
	switch err {
	case nil:
		c.JSON(200, result)
	case ErrInvalidToken:
		c.JSON(401, &errResult{err.Error()})
	case api.ErrLoginFailed:
		c.JSON(401, &errResult{err.Error()})
	default:
		c.JSON(500, &errResult{err.Error()})
	}
}
