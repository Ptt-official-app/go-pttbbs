package api

import (
	"github.com/gin-gonic/gin"
)

func OptionsWrapper(c *gin.Context) {
	processResult(c, nil, nil)
}
