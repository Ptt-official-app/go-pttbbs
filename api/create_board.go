package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const CREATE_BOARD_R = "/class/:cls/board"

type CreateBoardParams struct {
	Brdname      string            `json:"brdname" form:"brdname" url:"brdname"`
	BrdClass     []byte            `json:"class" form:"class" url:"class"`
	BrdTitle     []byte            `json:"title" form:"title" url:"title"`
	BMs          []bbs.UUserID     `json:"bms" form:"bms" url:"bms"`
	BrdAttr      ptttype.BrdAttr   `json:"brdattr" form:"brdattr" url:"brdattr"`
	Level        ptttype.PERM      `json:"level,omitempty" form:"level,omitempty" url:"level,omitempty"`
	ChessCountry ptttype.ChessCode `json:"chess_country,omitempty" form:"chess_country,omitempty" url:"chess_country,omitempty"`
	IsGroup      bool              `json:"is_group" form:"is_group" url:"is_group"`
}

type CreateBoardPath struct {
	ClsBid ptttype.Bid `uri:"cls" binding:"required"`
}

type CreateBoardResult struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

func CreateBoardWrapper(c *gin.Context) {
	params := &CreateBoardParams{}
	path := &CreateBoardPath{}
	LoginRequiredPathJSON(CreateBoard, params, path, c)
}

func CreateBoard(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*CreateBoardParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*CreateBoardPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	boardID, err := bbs.CreateBoard(
		uuserID,
		thePath.ClsBid,
		theParams.Brdname,
		theParams.BrdClass,
		theParams.BrdTitle,
		theParams.BMs,
		theParams.BrdAttr,
		theParams.Level,
		theParams.ChessCountry,
		theParams.IsGroup,
	)
	if err != nil {
		return nil, err
	}

	result = &CreateBoardResult{
		BBoardID: boardID,
	}

	return result, nil
}
