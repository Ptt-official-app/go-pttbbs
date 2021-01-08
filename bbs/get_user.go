package bbs

import "github.com/Ptt-official-app/go-pttbbs/ptt"

func GetUser(uuserID UUserID) (user *Userec, err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	userecRaw, err := ptt.GetUser(userIDRaw)
	if err != nil {
		return nil, err
	}

	userec2Raw, err := ptt.GetUser2(userIDRaw)
	if err != nil {
		return nil, err
	}

	user = NewUserecFromRaw(userecRaw, userec2Raw)

	return user, nil
}
