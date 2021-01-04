package bbs

import "github.com/Ptt-official-app/go-pttbbs/ptt"

func GetUser(uuserID UUserID) (user *Userec, err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	_, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, err
	}

	user = NewUserecFromRaw(userecRaw)

	return user, nil
}
