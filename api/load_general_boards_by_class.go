package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const LOAD_GENERAL_BOARDS_BY_CLASS_R = "/boards/byclass"

func LoadGeneralBoardsByClassWrapper(c *gin.Context) {
	params := &LoadGeneralBoardsParams{}
	LoginRequiredQuery(LoadGeneralBoardsByClass, params, c)
}

func LoadGeneralBoardsByClass(remoteAddr string, uuserID bbs.UUserID, params interface{}, c *gin.Context) (interface{}, error) {
	return loadGeneralBoardsCore(remoteAddr, uuserID, params, ptttype.BSORT_BY_CLASS)
}
