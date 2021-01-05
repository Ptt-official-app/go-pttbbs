package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func ChangePasswd(userID UUserID, origPasswd string, passwd string, ip string) (err error) {
	userIDRaw, err := userID.ToRaw()
	if err != nil {
		return err
	}

	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))

	return ptt.ChangePasswd(userIDRaw, []byte(origPasswd), []byte(passwd), ipRaw)
}
