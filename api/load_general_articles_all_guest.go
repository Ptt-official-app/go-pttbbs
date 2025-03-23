package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

func LoadGeneralArticlesAllGuestWrapper(c *gin.Context) {
	params := NewLoadGeneralArticlesParams()
	path := &LoadGeneralArticlesPath{}
	PathQuery(LoadGeneralArticlesAllGuest, params, path, c)
}

func LoadGeneralArticlesAllGuest(remoteAddr string, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*LoadGeneralArticlesParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*LoadGeneralArticlesPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	summary, nextIdx, nextCreateTime, isNewest, startNumIdx, err := bbs.LoadGeneralArticlesAllGuest(
		thePath.BBoardID,
		theParams.StartIdx,
		theParams.NArticles,
		theParams.Desc,
	)
	if err != nil {
		return nil, err
	}

	result = &LoadGeneralArticlesResult{
		Articles:       summary,
		IsNewest:       isNewest,
		NextIdx:        nextIdx,
		NextCreateTime: nextCreateTime,
		StartNumIdx:    startNumIdx,
	}

	return result, nil
}
