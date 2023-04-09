package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func WriteFavorites(uuserID UUserID, content []byte) (mtime types.Time4, err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return 0, ErrInvalidParams
	}

	return ptt.WriteFavorites(userIDRaw, content)
}
