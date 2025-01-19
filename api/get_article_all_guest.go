package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

func GetArticleAllGuestWrapper(c *gin.Context) {
	params := &GetArticleParams{}
	path := &GetArticlePath{}

	PathQuery(GetArticleAllGuest, params, path, c)
}

// GetArticle
//
// Require middleware to handle deleted files (by owner-name and by title).
// We will handle user-board relation and filename prefix.
//
// Require middleware to parse the content.
// Require middleware to take care of user-read-article.
func GetArticleAllGuest(remoteAddr string, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*GetArticleParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*GetArticlePath)
	if !ok {
		return nil, ErrInvalidPath
	}

	content, mtime, hash, err := bbs.GetArticleAllGuest(thePath.BBoardID, thePath.ArticleID, theParams.RetrieveTS, theParams.IsHash)
	if err != nil {
		return nil, err
	}

	result = &GetArticleResult{
		MTime:   mtime,
		Content: content,
		Hash:    hash,
	}

	return result, nil
}
