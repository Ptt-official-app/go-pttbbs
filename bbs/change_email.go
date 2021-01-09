package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func ChangeEmail(uuserID UUserID, email string) (err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return err
	}

	emailRaw := &ptttype.Email_t{}
	copy(emailRaw[:], []byte(email))

	return ptt.ChangeEmail(userIDRaw, emailRaw)
}
