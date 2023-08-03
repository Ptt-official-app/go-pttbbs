package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const CREATE_ARTICLE_R = "/board/:bid/article"

type CreateArticleParams struct {
	PostType []byte   `json:"class" form:"class" url:"class"`
	Title    []byte   `json:"title" form:"title" url:"title"`
	Content  [][]byte `json:"content" form:"content" url:"content"`
}

type CreateArticlePath struct {
	BBoardID bbs.BBoardID `uri:"bid" binding:"required"`
}

type CreateArticleResult *bbs.ArticleSummary

func CreateArticleWrapper(c *gin.Context) {
	params := &CreateArticleParams{}
	path := &CreateArticlePath{}
	LoginRequiredPathJSON(CreateArticle, params, path, c)
}

func CreateArticle(remoteAddr string, uuser bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*CreateArticleParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*CreateArticlePath)
	if !ok {
		return nil, ErrInvalidPath
	}

	summary, err := bbs.CreateArticle(uuser, thePath.BBoardID, theParams.PostType, theParams.Title, theParams.Content, remoteAddr)
	if err != nil {
		return nil, err
	}

	return CreateArticleResult(summary), nil
}
