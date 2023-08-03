package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const RELOAD_UHASH_R = "/admin/reloaduhash"

type ReloadUHashResult struct {
	Success bool `json:"success"`
}

func ReloadUHashWrapper(c *gin.Context) {
	LoginRequiredQuery(ReloadUHash, nil, c)
}

func ReloadUHash(remoteAddr string, uuserID bbs.UUserID, params interface{}, c *gin.Context) (result interface{}, err error) {
	err = bbs.ReloadUHash(uuserID)
	if err != nil {
		return nil, err
	}

	return &ReloadUHashResult{Success: true}, nil
}
