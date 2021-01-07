package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func CheckPasswd(uuserID UUserID, passwd string, ip string) (err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return ErrInvalidParams
	}

	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))

	return ptt.CheckPasswd(userIDRaw, []byte(passwd), ipRaw)
}
