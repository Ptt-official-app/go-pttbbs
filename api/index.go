package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const INDEX_R = "/"

type IndexParams struct{}

type IndexResult struct {
	Data string
}

func IndexWrapper(c *gin.Context) {
	params := &IndexParams{}
	LoginRequiredJSON(Index, params, c)
}

func Index(remoteAddr string, uuserID bbs.UUserID, params interface{}) (interface{}, error) {
	result := &IndexResult{Data: "index"}
	return result, nil
}
