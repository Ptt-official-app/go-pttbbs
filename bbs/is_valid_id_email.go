package bbs

import "github.com/Ptt-official-app/go-pttbbs/ptt"

func IsValidIDEmail(email string) (err error) {
	return ptt.CheckEmailAllowRejectLists(email)
}
