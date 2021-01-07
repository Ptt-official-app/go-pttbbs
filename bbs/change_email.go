package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func ChangeEmail(userID UUserID, email string) (err error) {
	userIDRaw, err := userID.ToRaw()
	if err != nil {
		return err
	}

	emailRaw := &ptttype.Email_t{}
	copy(emailRaw[:], []byte(email))

	return ptt.ChangeEmail(userIDRaw, emailRaw)
}
