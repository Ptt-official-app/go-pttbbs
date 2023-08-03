package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

const GET_FAV_R = "/user/:uid/favorites"

type GetFavoritesParams struct {
	RetrieveTS types.Time4 `json:"last_ts,omitempty" form:"last_ts,omitempty" url:"last_ts,omitempty"`
}

type GetFavoritesPath struct {
	UserID bbs.UUserID `uri:"uid" binding:"required"`
}

type GetFavoritesResult struct {
	MTime   types.Time4 `json:"mtime"`
	Content []byte      `json:"content"`
}

func GetFavoritesWrapper(c *gin.Context) {
	params := &GetFavoritesParams{}
	path := &GetFavoritesPath{}
	LoginRequiredPathQuery(GetFavorites, params, path, c)
}

func GetFavorites(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*GetFavoritesParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*GetFavoritesPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if uuserID != thePath.UserID {
		return nil, ErrInvalidUser
	}

	content, mtime, err := bbs.GetFavorites(thePath.UserID, theParams.RetrieveTS)
	if err != nil {
		return nil, err
	}

	result = &GetFavoritesResult{
		MTime:   mtime,
		Content: content,
	}

	return result, nil
}
