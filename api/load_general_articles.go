package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

const LOAD_GENERAL_ARTICLES_R = "/board/:bid/articles"

type LoadGeneralArticlesParams struct {
	StartIdx  string `json:"start_idx,omitempty" form:"start_idx,omitempty" url:"start_idx,omitempty"`
	NArticles int    `json:"max" form:"max" url:"max"`
	Desc      bool   `json:"desc,omitempty" form:"desc,omitempty" url:"desc"`
	IsSystem  bool   `json:"system,omitempty" form:"system,omitempty" url:"system"`
}

type LoadGeneralArticlesPath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

type LoadGeneralArticlesResult struct {
	Articles       []*bbs.ArticleSummary `json:"data"`
	IsNewest       bool                  `json:"is_newest"`
	NextIdx        string                `json:"next_idx"`
	NextCreateTime types.Time4           `json:"next_create_time"`
	StartNumIdx    ptttype.SortIdx       `json:"start_num_idx"`
}

func NewLoadGeneralArticlesParams() *LoadGeneralArticlesParams {
	return &LoadGeneralArticlesParams{
		Desc: true,
	}
}

func LoadGeneralArticlesWrapper(c *gin.Context) {
	params := NewLoadGeneralArticlesParams()
	path := &LoadGeneralArticlesPath{}
	LoginRequiredPathQuery(LoadGeneralArticles, params, path, c)
}

func LoadGeneralArticles(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*LoadGeneralArticlesParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*LoadGeneralArticlesPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if theParams.IsSystem {
		uuserID = bbs.UUserID(string(ptttype.STR_SYSOP))
	}

	summary, nextIdx, nextCreateTime, isNewest, startNumIdx, err := bbs.LoadGeneralArticles(
		uuserID,
		thePath.BBoardID,
		theParams.StartIdx,
		theParams.NArticles,
		theParams.Desc,
	)
	if err != nil {
		return nil, err
	}

	result = &LoadGeneralArticlesResult{
		Articles:       summary,
		IsNewest:       isNewest,
		NextIdx:        nextIdx,
		NextCreateTime: nextCreateTime,
		StartNumIdx:    startNumIdx,
	}

	return result, nil
}
