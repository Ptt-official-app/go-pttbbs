package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const CREATE_COMMENT_R = "/board/:bid/article/:aid/comment"

type CreateCommentParams struct {
	CommentType ptttype.CommentType `json:"type" form:"type" url:"type"`
	Content     []byte              `json:"content" form:"content" url:"content"`
}

type CreateCommentPath struct {
	BBoardID  bbs.BBoardID  `uri:"bid" binding:"required"`
	ArticleID bbs.ArticleID `uri:"aid" binding:"required"`
}

type CreateCommentResult struct {
	Content []byte      `json:"content"`
	MTime   types.Time4 `json:"mtime"`
}

func CreateCommentWrapper(c *gin.Context) {
	params := &CreateCommentParams{}
	path := &CreateCommentPath{}
	LoginRequiredPathJSON(CreateComment, params, path, c)
}

func CreateComment(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {

	theParams, ok := params.(*CreateCommentParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	if theParams.CommentType > ptttype.COMMENT_TYPE_BASIC {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*CreateCommentPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	logrus.Infof("CreateComment: commentType: %v", theParams.CommentType)

	content, mtime, err := bbs.CreateComment(uuserID, thePath.BBoardID, thePath.ArticleID, theParams.CommentType, theParams.Content, remoteAddr)
	if err != nil {
		return nil, err
	}

	result = &CreateCommentResult{Content: content, MTime: mtime}

	return result, nil
}
