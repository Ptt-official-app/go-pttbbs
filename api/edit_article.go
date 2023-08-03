package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

const EDIT_ARTICLE_R = "/board/:bid/article/:aid/edit"

type EditArticleParams struct {
	PostType []byte   `json:"class" form:"class" url:"class"`
	Title    []byte   `json:"title,omitempty" form:"title,omitempty" url:"title,omitempty"`
	Content  [][]byte `json:"content" form:"content" url:"content"`

	OldSZ  int           `json:"oldsz" form:"oldsz" url:"oldsz"`
	OldSum cmsys.Fnv64_t `json:"oldsum" form:"oldsum" url:"oldsum"`
}

type EditArticlePath struct {
	BBoardID  bbs.BBoardID  `uri:"bid" binding:"required"`
	ArticleID bbs.ArticleID `uri:"aid" binding:"required"`
}

type EditArticleResult struct {
	MTime     types.Time4 `json:"mtime"`
	Content   []byte      `json:"content"` // content contains all the necessary information.
	RealTitle []byte      `json:"title"`
	Class     []byte      `json:"class"`
	FullTitle []byte      `json:"full_title"`
}

func EditArticleWrapper(c *gin.Context) {
	params := &EditArticleParams{}
	path := &EditArticlePath{}
	LoginRequiredPathJSON(EditArticle, params, path, c)
}

func EditArticle(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*EditArticleParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*EditArticlePath)
	if !ok {
		return nil, ErrInvalidPath
	}

	content, mtime, theTitle, theClass, fullTitle, err := bbs.EditArticle(uuserID, thePath.BBoardID, thePath.ArticleID, theParams.PostType, theParams.Title, theParams.Content, theParams.OldSZ, theParams.OldSum, remoteAddr)
	if err != nil {
		return nil, err
	}

	r := &EditArticleResult{
		Content:   content,
		MTime:     mtime,
		RealTitle: theTitle,
		Class:     theClass,
		FullTitle: fullTitle,
	}

	return r, nil
}
