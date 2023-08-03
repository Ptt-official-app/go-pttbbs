package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const LOAD_CLASS_BOARDS_R = "/cls/:clsid/boards"

type LoadClassBoardsParams struct {
	BSortBy  ptttype.BSortBy `json:"sortby,omitempty" form:"sortby,omitempty" url:"sortby,omitempty"`
	IsSystem bool            `json:"system,omitempty" form:"system,omitempty" url:"system"`
}

type LoadClassBoardsPath struct {
	ClsID ptttype.Bid `uri:"clsid" binding:"required"`
}

type LoadClassBoardsResult struct {
	Boards []*bbs.BoardSummary `json:"data"`
}

func LoadClassBoardsWrapper(c *gin.Context) {
	params := &LoadClassBoardsParams{}
	path := &LoadClassBoardsPath{}
	LoginRequiredPathQuery(LoadClassBoards, params, path, c)
}

func LoadClassBoards(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (ret interface{}, err error) {
	theParams, ok := params.(*LoadClassBoardsParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*LoadClassBoardsPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if !theParams.IsSystem {
		return nil, ErrInvalidUser
	}

	uuserID = bbs.UUserID(string(ptttype.STR_SYSOP))

	summaries, err := bbs.LoadClassBoards(uuserID, thePath.ClsID, theParams.BSortBy)
	if err != nil {
		return nil, err
	}

	ret = &LoadClassBoardsResult{
		Boards: summaries,
	}

	return ret, nil
}
