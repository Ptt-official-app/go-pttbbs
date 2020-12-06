package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/names"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	log "github.com/sirupsen/logrus"
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
func LoginQuery(userID *ptttype.UserID_t, passwd []byte, ip *ptttype.IPv4_t) (*ptttype.UserecRaw, error) {
	if !names.IsValidUserID(userID) {
		return nil, ptttype.ErrInvalidUserID
	}

	_, cuser, err := initCurrentUser(userID)
	log.Debugf("after initCurrentUser: cuser: %v e: %v", cuser, err)
	if err != nil {
		return nil, err
	}

	isValid, err := cmbbs.CheckPasswd(cuser.PasswdHash[:], passwd)
	log.Debugf("mbbsd.LoginQuery: after CheckPasswd: isValid: %v e: %v", isValid, err)
	if err != nil {
		cmbbs.LogAttempt(userID, ip, true)
		return nil, err
	}

	if !isValid {
		cmbbs.LogAttempt(userID, ip, true)
		return nil, ptttype.ErrInvalidUserID
	}

	//XXX do post-user-login.

	return cuser, nil
}
