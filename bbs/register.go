package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func Register(
	username string,
	passwd string,
	ip string,
	email string,

	nickname []byte,
	realname []byte,
	career []byte,
	address []byte,
	over18 bool,
) (user *Userec, err error) {
	userIDRaw := &ptttype.UserID_t{}
	copy(userIDRaw[:], []byte(username))

	passwdRaw := []byte(passwd)

	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))

	emailRaw := &ptttype.Email_t{}
	copy(emailRaw[:], []byte(email))

	nicknameRaw := &ptttype.Nickname_t{}
	copy(nicknameRaw[:], nickname)

	realnameRaw := &ptttype.RealName_t{}
	copy(realnameRaw[:], realname)

	careerRaw := &ptttype.Career_t{}
	copy(careerRaw[:], []byte(career))

	addressRaw := &ptttype.Address_t{}
	copy(addressRaw[:], []byte(address))

	_, userRaw, err := ptt.NewRegister(
		userIDRaw,
		passwdRaw,
		ipRaw,
		emailRaw,
		false,
		false,

		nicknameRaw,
		realnameRaw,
		careerRaw,
		addressRaw,
		over18,
	)
	if err != nil {
		return nil, err
	}

	user = NewUserecFromRaw(userRaw)

	return user, nil
}
