package ptt

import (
	"errors"
	"fmt"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

var (
	ErrNotPermitted  = errors.New("not permitted")
	ErrInvalidParams = errors.New("invalid params")
)

func FatalLockedUser(userID *ptttype.UserID_t) error {
	return fmt.Errorf("[FATAL] System Error, Locked User! %v", types.CstrToString(userID[:]))
}
