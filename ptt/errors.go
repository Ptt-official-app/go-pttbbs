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
	ErrNewUtmp       = errors.New("unable to get new utmp")

	//talk.go
	ErrNoUser      = errors.New("no user")
	ErrTooManyMsgs = errors.New("too many msgs")
	ErrInvalidPID  = errors.New("invalid pid")
)

func FatalLockedUser(userID *ptttype.UserID_t) error {
	return fmt.Errorf("[FATAL] System Error, Locked User! %v", types.CstrToString(userID[:]))
}
