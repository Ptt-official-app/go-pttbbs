package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const LOAD_BOTTOM_ARTICLES_R = "/board/:bid/articles/bottom"

type LoadBottomArticlesPath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

type LoadBottomArticlesResult struct {
	Articles []*bbs.ArticleSummary `json:"data"`
}

func LoadBottomArticlesWrapper(c *gin.Context) {
	path := &LoadGeneralArticlesPath{}
	LoginRequiredPathQuery(LoadBottomArticles, nil, path, c)
}

func LoadBottomArticles(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	thePath, ok := path.(*LoadGeneralArticlesPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	summary, err := bbs.LoadBottomArticles(
		uuserID,
		thePath.BBoardID,
	)
	if err != nil {
		return nil, err
	}

	result = &LoadBottomArticlesResult{
		Articles: summary,
	}

	return result, nil
}
