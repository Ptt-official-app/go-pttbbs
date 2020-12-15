package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

const LOAD_GENERAL_ARTICLES_R = "/boards/:bid/articles"

type LoadGeneralArticlesParams struct {
	StartIdx  string `json:"start_idx,omitempty" form:"start_idx,omitempty" url:"start_idx,omitempty"`
	NArticles int    `json:"max,omitempty" form:"max,omitempty" url:"max,omitempty"`
}

type LoadGeneralArticlesPath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

type LoadGeneralArticlesResult struct {
	Articles []*bbs.ArticleSummary `json:"data"`
	NextIdx  string                `json:"next_idx"`
}

func LoadGeneralArticles(remoteAddr string, userID string, params interface{}, path interface{}) (result interface{}, err error) {
	loadGeneralArticlesParams, ok := params.(*LoadGeneralArticlesParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	loadGeneralArticlesPath, ok := path.(*LoadGeneralArticlesPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	summary, nextIdx, err := bbs.LoadGeneralArticles(userID, loadGeneralArticlesPath.BBoardID, loadGeneralArticlesParams.StartIdx, loadGeneralArticlesParams.NArticles)

	if err != nil {
		return nil, err
	}

	results := &LoadGeneralArticlesResult{
		Articles: summary,
		NextIdx:  nextIdx,
	}

	return results, nil
}
