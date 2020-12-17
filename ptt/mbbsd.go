package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

//LoginQuery
//
//Params
//	userID: userID
//	passwd: passwd
//	ip: ip
//
//Return
//	*UserecRaw: user
//  error: err
func LoginQuery(userID *ptttype.UserID_t, passwd []byte, ip *ptttype.IPv4_t) (uid ptttype.Uid, user *ptttype.UserecRaw, err error) {
	if !userID.IsValid() {
		return 0, nil, ptttype.ErrInvalidUserID
	}

	uid, user, err = InitCurrentUser(userID)
	if err != nil {
		return 0, nil, err
	}

	isValid, err := cmbbs.CheckPasswd(user.PasswdHash[:], passwd)
	if err != nil {
		cmbbs.LogAttempt(userID, ip, true)
		return 0, nil, err
	}

	if !isValid {
		cmbbs.LogAttempt(userID, ip, true)
		return 0, nil, ptttype.ErrInvalidUserID
	}

	//XXX do post-user-login.

	return uid, user, nil
}
