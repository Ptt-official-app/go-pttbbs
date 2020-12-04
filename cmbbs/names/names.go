package names

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//IsValidUserID
//
//Params
//	userID: user-id
//
//Return
//	bool: is valid user-id
func IsValidUserID(userID *ptttype.UserID_t) bool {
	if userID == nil {
		return false
	}

	theLen := types.Cstrlen(userID[:])
	if theLen < 2 || theLen > ptttype.IDLEN {
		return false
	}

	if !types.Isalpha(userID[0]) {
		return false
	}

	for idx, c := range userID {
		if idx == theLen {
			break
		}

		if !types.Isalnum(c) {
			return false
		}
	}

	return true
}
