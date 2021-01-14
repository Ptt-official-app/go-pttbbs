package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const LOAD_GENERAL_ARTICLES_R = "/board/:bid/articles"

type LoadGeneralArticlesParams struct {
	StartIdx  string `json:"start_idx,omitempty" form:"start_idx,omitempty" url:"start_idx,omitempty"`
	NArticles int    `json:"max,omitempty" form:"max,omitempty" url:"max,omitempty"`
}

type LoadGeneralArticlesPath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

type LoadGeneralArticlesResult struct {
	Articles []*bbs.ArticleSummary `json:"data"`
	IsNewest bool                  `json:"is_newest"`
	NextIdx  string                `json:"next_idx"`
}

func LoadGeneralArticlesWrapper(c *gin.Context) {
	params := &LoadGeneralArticlesParams{}
	path := &LoadGeneralArticlesPath{}
	LoginRequiredPathQuery(LoadGeneralArticles, params, path, c)
}

func LoadGeneralArticles(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	loadGeneralArticlesParams, ok := params.(*LoadGeneralArticlesParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	loadGeneralArticlesPath, ok := path.(*LoadGeneralArticlesPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	summary, nextIdx, isNewest, err := bbs.LoadGeneralArticles(uuserID, loadGeneralArticlesPath.BBoardID, loadGeneralArticlesParams.StartIdx, loadGeneralArticlesParams.NArticles)

	if err != nil {
		return nil, err
	}

	result = &LoadGeneralArticlesResult{
		Articles: summary,
		IsNewest: isNewest,
		NextIdx:  nextIdx,
	}

	return result, nil
}
