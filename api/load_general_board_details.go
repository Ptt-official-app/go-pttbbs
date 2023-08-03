package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const LOAD_GENERAL_BOARD_DETAILS_R = "/boards/detail"

type LoadGeneralBoardDetailsParams struct {
	StartIdx string `json:"start_idx,omitempty" form:"start_idx,omitempty" url:"start_idx,omitempty"`
	NBoards  int    `json:"max" form:"max" url:"max"`
	Asc      bool   `json:"asc,omitempty" form:"asc,omitempty" url:"asc"`
	IsSystem bool   `json:"system,omitempty" form:"system,omitempty" url:"system"`
}

type LoadGeneralBoardDetailsResult struct {
	Boards  []*bbs.BoardDetail `json:"data"`
	NextIdx string             `json:"next_idx"`
}

func NewLoadGeneralBoardDetailsParams() *LoadGeneralBoardDetailsParams {
	return &LoadGeneralBoardDetailsParams{
		Asc: true,
	}
}

func LoadGeneralBoardDetailsWrapper(c *gin.Context) {
	params := NewLoadGeneralBoardDetailsParams()
	LoginRequiredQuery(LoadGeneralBoardDetails, params, c)
}

func LoadGeneralBoardDetails(remoteAddr string, uuserID bbs.UUserID, params interface{}, c *gin.Context) (interface{}, error) {
	return loadGeneralBoardDetailsCore(remoteAddr, uuserID, params, ptttype.BSORT_BY_NAME)
}

func loadGeneralBoardDetailsCore(remoteAddr string, uuserID bbs.UUserID, params interface{}, bsortBy ptttype.BSortBy) (interface{}, error) {
	theParams, ok := params.(*LoadGeneralBoardDetailsParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	if theParams.IsSystem {
		uuserID = bbs.UUserID(string(ptttype.STR_SYSOP))
	}

	if uuserID != bbs.UUserID(string(ptttype.STR_SYSOP)) {
		return nil, ErrInvalidUser
	}

	details, nextIdx, err := bbs.LoadGeneralBoardDetails(uuserID, theParams.StartIdx, theParams.NBoards, theParams.Asc, bsortBy)
	if err != nil {
		return nil, err
	}

	results := &LoadGeneralBoardDetailsResult{
		Boards:  details,
		NextIdx: nextIdx,
	}

	return results, nil
}
