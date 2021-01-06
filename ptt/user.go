package ptt

import (
	"os"
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"

	log "github.com/sirupsen/logrus"
)

//killUser
//
//Assume correct uid / userID correspondance.
func killUser(uid ptttype.Uid, userID *ptttype.UserID_t) error {
	if uid <= 0 || userID == nil {
		return ptttype.ErrInvalidUserID
	}

	err := friendDeleteAll(userID, ptttype.FRIEND_ALOHA)
	if err != nil {
		log.Errorf("killUser: unable to friend-delete-all: uid: %v e: %v", uid, err)
	}

	err = tryDeleteHomePath(userID)
	if err != nil {
		log.Errorf("killUser: unable to delete home-path: userID: %v e: %v", userID, err)
	}

	emptyUser := &ptttype.UserecRaw{}
	err = passwdSyncUpdate(uid, emptyUser)
	if err != nil {
		log.Errorf("killUser: unable to passwd-sync-update emptyUser: uid: %v e: %v", uid, err)
	}

	return nil
}

func tryDeleteHomePath(userID *ptttype.UserID_t) error {
	homePath := path.SetHomePath(userID)
	dstPath := strings.Join([]string{ptttype.BBSHOME, ptttype.DIR_TMP, types.CstrToString(userID[:])}, string(os.PathSeparator))

	if !types.IsDir(homePath) {
		return nil
	}

	if err := types.Rename(homePath, dstPath); err != nil {
		return err
	}

	if err := os.RemoveAll(homePath); err != nil {
		return err
	}

	return nil
}

//https://github.com/ptt/pttbbs/blob/master/mbbsd/user.c#L961
//https://github.com/ptt/pttbbs/blob/master/mbbsd/user.c#L1194
func ChangePasswd(userID *ptttype.UserID_t, origPasswd []byte, passwd []byte, ip *ptttype.IPv4_t) (err error) {
	if userID == nil || userID[0] == 0 {
		return ptttype.ErrInvalidUserID
	}

	uid, err := cache.SearchUserRaw(userID, nil)
	if err != nil {
		return err
	}

	userPasswdHash, err := cmbbs.PasswdQueryPasswd(uid)
	if err != nil {
		return err
	}

	isValid, err := cmbbs.CheckPasswd(userPasswdHash[:], origPasswd)
	if err != nil {
		cmbbs.LogAttempt(userID, ip, true)
		return err
	}

	if !isValid {
		cmbbs.LogAttempt(userID, ip, true)
		return ptttype.ErrInvalidUserID
	}

	genNewPasswdHash, err := cmbbs.GenPasswd(passwd)
	if err != nil {
		return ErrInvalidParams
	}

	err = cmbbs.PasswdUpdatePasswd(uid, genNewPasswdHash)
	if err != nil {
		return err
	}

	return nil
}
