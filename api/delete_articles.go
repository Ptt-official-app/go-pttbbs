package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const DELETE_ARTICLES_R = "/board/:bid/deletearticles"

type DeleteArticlesParams struct {
	ArticleIDs []bbs.ArticleID `json:"aids" form:"aids" url:"aids" binding:"required"`
}

type DeleteArticlesPath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

type DeleteArticlesResult struct{}

func DeleteArticlesWrapper(c *gin.Context) {
	params := &DeleteArticlesParams{}
	path := &DeleteArticlesPath{}
	LoginRequiredPathJSON(DeleteArticles, params, path, c)
}

func DeleteArticles(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*DeleteArticlesParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*DeleteArticlesPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	err = bbs.DeleteArticles(uuserID, thePath.BBoardID, theParams.ArticleIDs, remoteAddr)
	if err != nil {
		return nil, err
	}
	result = DeleteArticlesResult{}
	return result, nil
}
