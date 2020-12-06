package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//verifyDBSet
//
//XXX to implement
func verifyDBSet(userID *ptttype.UserID_t, generation int64, vmethod ptttype.VerifyDBVMethod, vkey []byte, timestamp int64) error {
	return types.ErrNotImplemented
}

//verifyDBCountByVerify
//
//XXX to implement
func verifyDBCountByVerify(vmethod ptttype.VerifyDBVMethod, vkey []byte) (countSelf int, countOther int, err error) {
	return 0, 0, types.ErrNotImplemented
}
