package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const SET_USER_PERM_R = "/admin/user/:uid/setperm"

type SetUserPermParams struct {
	Perm ptttype.PERM `json:"perm"`
}

type SetUserPermPath struct {
	UserID bbs.UUserID `uri:"uid" binding:"required"`
}

type SetUserPermResult struct {
	Perm ptttype.PERM `json:"perm"`
}

func SetUserPermWrapper(c *gin.Context) {
	params := &SetUserPermParams{}
	path := &SetUserPermPath{}
	LoginRequiredPathJSON(SetUserPerm, params, path, c)
}

func SetUserPerm(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error) {
	theParams, ok := params.(*SetUserPermParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	thePath, ok := path.(*SetUserPermPath)
	if !ok {
		return nil, ErrInvalidPath
	}

	perm, err := bbs.SetUserPerm(uuserID, thePath.UserID, theParams.Perm)
	if err != nil {
		return nil, err
	}

	return &SetUserPermResult{Perm: perm}, nil
}
