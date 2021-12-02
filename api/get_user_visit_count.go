package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const GET_USER_VISIT_COUNT_R = "/uservisitcount"

type GetUserVisitCountResult struct {
	Total int32 `json:"total"`
}

func GetUserVisitCountWrapper(c *gin.Context) {
	Query(GetUserVisitCount, nil, c)
}

func GetUserVisitCount(remoteAddr string, params interface{}) (interface{}, error) {
	total := bbs.GetUserVisitCount()
	return &GetUserVisitCountResult{total}, nil
}
