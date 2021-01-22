package bbs

import "github.com/Ptt-official-app/go-pttbbs/ptt"

func CheckExistsUser(username string) (uuserID UUserID, err error) {
	uuserID = UUserID(username)
	uuserIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return "", ErrInvalidParams
	}

	uid, err := ptt.GetUid(uuserIDRaw)
	if err != nil {
		return "", err
	}
	if !uid.IsValid() {
		return "", nil
	}
	return uuserID, nil
}
