package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

//Login
//
//XXX need to check for the permission
func Login(username string, passwd string, ip string) (*Userec, error) {
	userIDRaw := &ptttype.UserID_t{}
	copy(userIDRaw[:], []byte(username))
	passwdRaw := []byte(passwd)
	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))

	_, userRaw, err := ptt.Login(userIDRaw, passwdRaw, ipRaw)
	if err != nil {
		return nil, err
	}

	user := NewUserecFromRaw(userRaw)

	return user, nil
}
