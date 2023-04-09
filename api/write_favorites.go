package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

const WRITE_FAV_R = "/user/:uid/favorites/post"

type WriteFavoritesParams struct {
	Content []byte `json:"content,omitempty" form:"content,omitempty" url:"content,omitempty"`
}

type WriteFavoritesPath struct {
	UserID bbs.UUserID `uri:"uid" binding:"required"`
}

type WriteFavoritesResult struct {
	MTime types.Time4 `json:"mtime"`
}

func WriteFavoritesWrapper(c *gin.Context) {
	params := &GetFavoritesParams{}
	path := &GetFavoritesPath{}
	LoginRequiredPathQuery(WriteFavorites, params, path, c)
}

func WriteFavorites(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}) (result interface{}, err error) {
	theParams, ok := params.(*WriteFavoritesParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*WriteFavoritesPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	if uuserID != thePath.UserID {
		return nil, ErrInvalidUser
	}

	mtime, err := bbs.WriteFavorites(thePath.UserID, theParams.Content)
	if err != nil {
		return nil, err
	}

	result = &WriteFavoritesResult{
		MTime: mtime,
	}

	return result, nil
}
