package ptt

import (
	"fmt"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

var ()

func FatalLockedUser(userID *ptttype.UserID_t) error {
	return fmt.Errorf("[FATAL] System Error, Locked User! %v", types.CstrToString(userID[:]))
}
