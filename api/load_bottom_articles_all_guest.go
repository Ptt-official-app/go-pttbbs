package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

func LoadBottomArticlesAllGuestWrapper(c *gin.Context) {
	path := &LoadGeneralArticlesPath{}
	PathQuery(LoadBottomArticlesAllGuest, nil, path, c)
}

func LoadBottomArticlesAllGuest(remoteAddr string, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	thePath, ok := path.(*LoadGeneralArticlesPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	summary, err := bbs.LoadBottomArticlesAllGuest(
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
