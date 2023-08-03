package api

import (
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

const GET_VERSION_R = "/version"

type GetVersionResult struct {
	Version    string `json:"version"`
	GitVersion string `json:"commit"`
}

func GetVersionWrapper(c *gin.Context) {
	Query(GetVersion, nil, c)
}

func GetVersion(remoteAddr string, params interface{}, c *gin.Context) (interface{}, error) {
	return &GetVersionResult{
		Version:    types.VERSION,
		GitVersion: types.GIT_VERSION,
	}, nil
}
