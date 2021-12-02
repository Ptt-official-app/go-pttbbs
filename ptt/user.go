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

func GetUID(userID *ptttype.UserID_t) (uid ptttype.UID, err error) {
	return cache.SearchUserRaw(userID, nil)
}

//killUser
//
//Assume correct uid / userID correspondance.
func killUser(uid ptttype.UID, userID *ptttype.UserID_t) error {
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

	if err := os.RemoveAll(homePath); err != nil { // nolint
		return err
	}

	return nil
}

// https://github.com/ptt/pttbbs/blob/master/mbbsd/user.c#L961
// https://github.com/ptt/pttbbs/blob/master/mbbsd/user.c#L1194
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

func CheckPasswd(userID *ptttype.UserID_t, passwd []byte, ip *ptttype.IPv4_t) (err error) {
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

	isValid, err := cmbbs.CheckPasswd(userPasswdHash[:], passwd)
	if err != nil {
		cmbbs.LogAttempt(userID, ip, true)
		return err
	}

	if !isValid {
		cmbbs.LogAttempt(userID, ip, true)
		return ptttype.ErrInvalidUserID
	}

	return nil
}

func ChangeEmail(userID *ptttype.UserID_t, email *ptttype.Email_t) (err error) {
	if userID == nil || userID[0] == 0 {
		return ptttype.ErrInvalidUserID
	}

	uid, err := cache.SearchUserRaw(userID, nil)
	if err != nil {
		return err
	}

	err = cmbbs.PasswdUpdateEmail(uid, email)
	if err != nil {
		return err
	}

	return nil
}

func ChangeUserLevel2(userID *ptttype.UserID_t, perm ptttype.PERM2, isSet bool) (userLevel2 ptttype.PERM2, err error) {
	if userID == nil || userID[0] == 0 {
		return ptttype.PERM2_INVALID, ptttype.ErrInvalidUserID
	}

	err = cmbbs.PasswdUpdateUserLevel2(userID, perm, isSet)
	if err != nil {
		return ptttype.PERM2_INVALID, err
	}

	return cmbbs.PasswdGetUserLevel2(userID)
}

// SetUserPerm
// https://github.com/ptt/pttbbs/blob/master/mbbsd/user.c#L1166
func SetUserPerm(userec *ptttype.UserecRaw, setUID ptttype.UID, setUserec *ptttype.UserecRaw, perm ptttype.PERM) (newPerm ptttype.PERM, err error) {
	setUserec.UserLevel = perm

	err = passwdSyncUpdate(setUID, setUserec)
	if err != nil {
		return ptttype.PERM_INVALID, err
	}

	return perm, nil
}

func GetUserVisitCount() int32 {
	return cache.GetUTMPNumber()
}
