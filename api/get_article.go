package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

const GET_ARTICLE_R = "/board/:bid/article/:aid"

type GetArticleParams struct {
	RetrieveTS types.Time4 `json:"last_ts,omitempty" form:"last_ts,omitempty" url:"last_ts,omitempty"`
}

type GetArticlePath struct {
	BBoardID  bbs.BBoardID  `uri:"bid" binding:"required"`
	ArticleID bbs.ArticleID `uri:"aid" binding:"required"`
}

type GetArticleResult struct {
	MTime   types.Time4 `json:"mtime"`
	Content []byte      `json:"content"` // content contains all the necessary information.
}

func GetArticleWrapper(c *gin.Context) {
	params := &GetArticleParams{}
	path := &GetArticlePath{}

	LoginRequiredPathQuery(GetArticle, params, path, c)
}

//GetArticle
//
//Require middleware to handle deleted files (by owner-name and by title).
//We will handle user-board relation and filename prefix.
//
//Require middleware to parse the content.
//Require middleware to take care of user-read-article.
func GetArticle(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*GetArticleParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*GetArticlePath)
	if !ok {
		return nil, ErrInvalidPath
	}

	content, mtime, err := bbs.GetArticle(uuserID, thePath.BBoardID, thePath.ArticleID, theParams.RetrieveTS)
	if err != nil {
		return nil, err
	}

	result = &GetArticleResult{
		MTime:   mtime,
		Content: content,
	}

	return result, nil
}
