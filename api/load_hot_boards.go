package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const LOAD_HOT_BOARDS_R = "/boards/popular"

type LoadHotBoardsResult struct {
	Boards []*bbs.BoardSummary `json:"data"`
}

func LoadHotBoardsWrapper(c *gin.Context) {
	LoginRequiredQuery(LoadHotBoards, nil, c)
}

// We have only 128 hot-boards.
func LoadHotBoards(remoteAddr string, uuserID bbs.UUserID, params interface{}, c *gin.Context) (result interface{}, err error) {
	summary, err := bbs.LoadHotBoards(uuserID)
	if err != nil {
		return nil, err
	}

	result = &LoadHotBoardsResult{
		Boards: summary,
	}

	return result, nil
}
