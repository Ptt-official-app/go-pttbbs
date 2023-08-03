package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

const CROSS_POST_R = "/board/:bid/article/:aid/crosspost"

type CrossPostParams struct {
	XBoardID bbs.BBoardID `json:"xbid" form:"xbid" url:"xbid"`
}

type CrossPostPath struct {
	BBoardID  bbs.BBoardID  `uri:"bid" binding:"required"`
	ArticleID bbs.ArticleID `uri:"aid" binding:"required"`
}

func CrossPostWrapper(c *gin.Context) {
	params := &CrossPostParams{}
	path := &CrossPostPath{}
	LoginRequiredPathJSON(CrossPost, params, path, c)
}

type CrossPostResult struct {
	// bbs.ArticleSummary
	ArticleSummary *bbs.ArticleSummary `json:"article"`

	Comment      []byte      `json:"comment"`
	CommentMTime types.Time4 `json:"comment_mtime"`
}

func CrossPost(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*CrossPostParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*CrossPostPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	articleSummary, comment, commentMTime, err := bbs.CrossPost(uuserID, thePath.BBoardID, thePath.ArticleID, theParams.XBoardID, remoteAddr)
	if err != nil {
		return nil, err
	}

	result = &CrossPostResult{ArticleSummary: articleSummary, Comment: comment, CommentMTime: commentMTime}

	return result, nil
}
