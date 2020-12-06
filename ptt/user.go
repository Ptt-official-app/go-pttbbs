package ptt

import (
	"os"
	"strings"

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

func hasUserPerm(user *ptttype.UserecRaw, perm ptttype.PERM) bool {
	return user.UserLevel&perm != 0
}
