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
	ErrNoUser        = errors.New("no user")
	ErrTooManyMsgs   = errors.New("too many msgs")
	ErrInvalidPID    = errors.New("invalid pid")
	ErrInvalidEmail  = errors.New("invalid email")
	ErrTooManyBoards = errors.New("too many boards")
	ErrCooldown      = errors.New("cooldown")
	ErrReadOnly      = errors.New("read only")
	ErrBanned        = errors.New("banned")
	ErrNoPost        = errors.New("no post")
	ErrRestricted    = errors.New("restricted")
	ErrViolateLaw    = errors.New("violate law")
)

func FatalLockedUser(userID *ptttype.UserID_t) error {
	return fmt.Errorf("[FATAL] System Error, Locked User! %v", types.CstrToString(userID[:]))
}
