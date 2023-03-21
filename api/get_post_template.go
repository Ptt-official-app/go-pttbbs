package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

const GET_POST_TEMPLATE_R = "/board/:bid/posttemplate/:tid"

type GetPostTemplateParams struct {
	RetrieveTS types.Time4 `json:"last_ts,omitempty" form:"last_ts,omitempty" url:"last_ts,omitempty"`
	IsSystem   bool        `json:"system,omitempty" form:"system,omitempty" url:"system,omitempty"`
	IsHash     bool        `json:"hash,omitempty" form:"hash,omitempty" url:"hash,omitempty"`
}

type GetPostTemplatePath struct {
	BBoardID   bbs.BBoardID    `uri:"bid" binding:"required"`
	TemplateID ptttype.SortIdx `uri:"tid" binding:"required"`
}

type GetPostTemplateResult struct {
	MTime   types.Time4   `json:"mtime"`
	Content []byte        `json:"content"` // content contains all the necessary information.
	Hash    cmsys.Fnv64_t `json:"hash"`
}

func GetPostTemplateWrapper(c *gin.Context) {
	params := &GetPostTemplateParams{}
	path := &GetPostTemplatePath{}

	LoginRequiredPathQuery(GetPostTemplate, params, path, c)
}

func GetPostTemplate(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*GetPostTemplateParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*GetPostTemplatePath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if theParams.IsSystem {
		uuserID = bbs.UUserID(string(ptttype.STR_SYSOP))
	}

	content, mtime, hash, err := bbs.GetPostTemplate(uuserID, thePath.BBoardID, thePath.TemplateID, theParams.RetrieveTS, theParams.IsHash)
	if err != nil {
		return nil, err
	}

	result = &GetPostTemplateResult{
		MTime:   mtime,
		Content: content,
		Hash:    hash,
	}

	return result, nil
}
