package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func GetFavorites(uuserID UUserID, retrieveTS types.Time4) (content []byte, mtime types.Time4, err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, 0, ErrInvalidParams
	}

	return ptt.GetFavorites(userIDRaw, retrieveTS)
}
