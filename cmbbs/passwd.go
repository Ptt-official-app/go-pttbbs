package cmbbs

import "C"
import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"os"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/crypt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/sem"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

// GenPasswd
//
// If passwd as empty: return empty passwd (unable to login)
func GenPasswd(passwd []byte) (passwdHash *ptttype.Passwd_t, err error) {
	if passwd[0] == 0 {
		return &ptttype.Passwd_t{}, nil
	}

	num := rand.Intn(65536)
	saltc := [2]byte{
		byte(num & 0x7f),
		byte((num >> 8) & 0x7f),
	}

	result, err := crypt.Fcrypt(passwd, saltc[:])
	passwdHash = &ptttype.Passwd_t{}
	copy(passwdHash[:], result[:])
	return passwdHash, err
}

// CheckPasswd
//
// Params
//
//	expected: expected-passwd-hash
//	input: input-passwd
//
// Return
//
//	bool: true: good (password matched). false: bad (password not matched).
//	error: err
func CheckPasswd(expected []byte, input []byte) (bool, error) {
	pw, err := crypt.Fcrypt(input, expected)
	if err != nil {
		return false, err
	}
	return bytes.Equal(pw, expected), nil // requires the passwd-hash be exact match.
}

func LogAttempt(userID *ptttype.UserID_t, ip *ptttype.IPv4_t, isWithUserHome bool) {
}

// PasswdLoadUser
//
// Params
//
//	userID: user-id
//
// Return
//
//	Uid: uid
//	*ptttype.UserecRaw: user.
//	error: err.
func PasswdLoadUser(userID *ptttype.UserID_t) (uid ptttype.UID, user *ptttype.UserecRaw, err error) {
	if types.IS_ALL_GUEST {
		return ptttype.GUEST_UID, ptttype.GUEST, nil
	}

	if userID == nil || userID[0] == 0 {
		return 0, nil, ptttype.ErrInvalidUserID
	}

	uid, err = cache.SearchUserRaw(userID, nil)
	if err != nil {
		return 0, nil, err
	}

	if !uid.IsValid() {
		return 0, nil, ptttype.ErrInvalidUserID
	}

	user, err = PasswdQuery(uid)
	if err != nil {
		return 0, nil, err
	}

	return uid, user, nil
}

// PasswdQuery
//
// Params
//
//	uid: uid
//
// Return
//
//	*ptttype.UserecRaw: user.
//	error: err.
func PasswdQuery(uid ptttype.UID) (*ptttype.UserecRaw, error) {
	if types.IS_ALL_GUEST {
		return ptttype.GUEST, nil
	}

	if !uid.IsValid() {
		return nil, ptttype.ErrInvalidUserID
	}

	file, err := os.Open(ptttype.FN_PASSWD)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	user := &ptttype.UserecRaw{}
	uidInFile := uid.ToUIDInStore()
	offset := int64(ptttype.USEREC_RAW_SZ) * int64(uidInFile)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}
	err = types.BinaryRead(file, binary.LittleEndian, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// PasswdQueryPasswd
//
// Params
//
//	uid: uid
//
// Return
//
//	*ptttype.UserecRaw: user.
//	error: err.
func PasswdQueryPasswd(uid ptttype.UID) (passwdHash *ptttype.Passwd_t, err error) {
	if types.IS_ALL_GUEST {
		return &ptttype.GUEST.PasswdHash, nil
	}

	if !uid.IsValid() {
		return nil, ptttype.ErrInvalidUserID
	}

	file, err := os.Open(ptttype.FN_PASSWD)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	passwdHash = &ptttype.Passwd_t{}
	uidInFile := uid.ToUIDInStore()
	offset := int64(ptttype.USEREC_RAW_SZ)*int64(uidInFile) + int64(unsafe.Offsetof(ptttype.USEREC_RAW.PasswdHash))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}
	err = types.BinaryRead(file, binary.LittleEndian, passwdHash)
	if err != nil {
		return nil, err
	}

	return passwdHash, nil
}

// PasswdQueryUserLevel
//
// Params
//
//	uid: uid
//
// Return
//
//	userLevel: userLevel.
//	error: err.
func PasswdQueryUserLevel(uid ptttype.UID) (userLevel ptttype.PERM, err error) {
	if types.IS_ALL_GUEST {
		return ptttype.GUEST.UserLevel, nil
	}

	if !uid.IsValid() {
		return ptttype.PERM_INVALID, ptttype.ErrInvalidUserID
	}

	file, err := os.Open(ptttype.FN_PASSWD)
	if err != nil {
		return ptttype.PERM_INVALID, err
	}
	defer file.Close()

	uidInFile := uid.ToUIDInStore()
	offset := int64(ptttype.USEREC_RAW_SZ)*int64(uidInFile) + int64(unsafe.Offsetof(ptttype.USEREC_RAW.UserLevel))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return ptttype.PERM_INVALID, err
	}
	err = types.BinaryRead(file, binary.LittleEndian, &userLevel)
	if err != nil {
		return ptttype.PERM_INVALID, err
	}

	return userLevel, nil
}

func PasswdUpdate(uid ptttype.UID, user *ptttype.UserecRaw) error {
	if types.IS_ALL_GUEST {
		return nil
	}

	if !uid.IsValid() {
		return cache.ErrInvalidUID
	}

	file, err := os.OpenFile(ptttype.FN_PASSWD, os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	uidInFile := uid.ToUIDInStore()
	_, err = file.Seek(int64(ptttype.USEREC_RAW_SZ)*int64(uidInFile), 0)
	if err != nil {
		return err
	}

	err = types.BinaryWrite(file, binary.LittleEndian, user)
	if err != nil {
		return err
	}

	return nil
}

func PasswdUpdatePasswd(uid ptttype.UID, passwdHash *ptttype.Passwd_t) error {
	if types.IS_ALL_GUEST {
		return nil
	}

	if !uid.IsValid() {
		return cache.ErrInvalidUID
	}

	file, err := os.OpenFile(ptttype.FN_PASSWD, os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	uidInFile := uid.ToUIDInStore()
	offset := int64(ptttype.USEREC_RAW_SZ)*int64(uidInFile) + int64(unsafe.Offsetof(ptttype.USEREC_RAW.PasswdHash))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	err = types.BinaryWrite(file, binary.LittleEndian, passwdHash)
	if err != nil {
		return err
	}

	return nil
}

func PasswdUpdateEmail(uid ptttype.UID, email *ptttype.Email_t) error {
	if types.IS_ALL_GUEST {
		return nil
	}

	if !uid.IsValid() {
		return cache.ErrInvalidUID
	}

	file, err := os.OpenFile(ptttype.FN_PASSWD, os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	uidInFile := uid.ToUIDInStore()
	offset := int64(ptttype.USEREC_RAW_SZ)*int64(uidInFile) + int64(unsafe.Offsetof(ptttype.USEREC_RAW.Email))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	err = types.BinaryWrite(file, binary.LittleEndian, email)
	if err != nil {
		return err
	}

	return nil
}

func PasswdGetUser2(userID *ptttype.UserID_t) (user *ptttype.Userec2Raw, err error) {
	if types.IS_ALL_GUEST {
		return ptttype.GUEST2, nil
	}

	filename, err := path.SetHomeFile(userID, ptttype.FN_PASSWD2)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &ptttype.Userec2Raw{}, nil
		}

		return nil, err
	}
	defer file.Close()

	user = &ptttype.Userec2Raw{}

	err = types.BinaryRead(file, binary.LittleEndian, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func PasswdGetUserLevel2(userID *ptttype.UserID_t) (userLevel2 ptttype.PERM2, err error) {
	if types.IS_ALL_GUEST {
		return ptttype.GUEST2.UserLevel2, nil
	}

	filename, err := path.SetHomeFile(userID, ptttype.FN_PASSWD2)
	if err != nil {
		return ptttype.PERM2_INVALID, err
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return ptttype.PERM2_INVALID, err
	}
	defer file.Close()

	offset := int64(unsafe.Offsetof(ptttype.USEREC2_RAW.UserLevel2))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return ptttype.PERM2_INVALID, err
	}

	err = types.BinaryRead(file, binary.LittleEndian, &userLevel2)
	if err != nil {
		return ptttype.PERM2_INVALID, err
	}

	return userLevel2, nil
}

func PasswdUpdateUserLevel2(userID *ptttype.UserID_t, perm ptttype.PERM2, isSet bool) (err error) {
	if types.IS_ALL_GUEST {
		return nil
	}

	filename, err := path.SetHomeFile(userID, ptttype.FN_PASSWD2)
	if err != nil {
		return err
	}

	err = passwdCheckPasswd2(filename)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	offset := int64(unsafe.Offsetof(ptttype.USEREC2_RAW.UserLevel2))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	userLevel2 := ptttype.PERM2(0)
	err = types.BinaryRead(file, binary.LittleEndian, &userLevel2)
	if isSet {
		userLevel2 |= perm
	} else {
		userLevel2 &= ^perm
	}

	offset = int64(unsafe.Offsetof(ptttype.USEREC2_RAW.UserLevel2))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}
	err = types.BinaryWrite(file, binary.LittleEndian, &userLevel2)
	if err != nil {
		return err
	}

	offset = int64(unsafe.Offsetof(ptttype.USEREC2_RAW.UpdateTS))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	updateTS := types.NowTS()
	err = types.BinaryWrite(file, binary.LittleEndian, &updateTS)
	if err != nil {
		return err
	}

	return nil
}

func passwdCheckPasswd2(filename string) (err error) {
	if types.IS_ALL_GUEST {
		return nil
	}

	stat, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			return err
		}
		defer file.Close()

		userec2 := &ptttype.Userec2Raw{
			Version: ptttype.PASSWD2_VERSION,
		}

		return types.BinaryWrite(file, binary.LittleEndian, userec2)
	}

	diffSize := int64(ptttype.USEREC2_RAW_SZ) - stat.Size()
	if diffSize == 0 {
		return nil
	}
	if diffSize < 0 {
		return ErrInvalidPasswd2Size
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	theBytes := make([]byte, diffSize)
	_, err = file.Write(theBytes)
	if err != nil {
		return err
	}
	return nil
}

func PasswdInit() error {
	if types.IS_ALL_GUEST {
		return nil
	}

	if Sem != nil {
		if Sem.SemID == 0 {
			log.Errorf("PasswdInit: SemID is invalid")
			return ErrSemInvalid
		}
		return nil
	}

	var err error
	Sem, err = sem.SemGet(ptttype.PASSWDSEM_KEY, 1, sem.SEM_R|sem.SEM_A|sem.IPC_CREAT|sem.IPC_EXCL)
	if err != nil {
		if os.IsExist(err) {
			Sem, err = sem.SemGet(ptttype.PASSWDSEM_KEY, 1, sem.SEM_R|sem.SEM_A)
			if err != nil {
				return err
			}

			return nil
		} else {
			return err
		}
	}

	err = Sem.SetVal(0, 1)
	if err != nil {
		return err
	}

	_, err = Sem.GetVal(0)
	if err != nil {
		return err
	}

	return nil
}

func PasswdLock() error {
	if types.IS_ALL_GUEST {
		return nil
	}

	return Sem.Wait(0)
}

func PasswdUnlock() error {
	if types.IS_ALL_GUEST {
		return nil
	}

	return Sem.Post(0)
}

// PasswdDestroy
//
// XXX [WARNING] know what you are doing before using Close!.
// This is to be able to close the semaphore for the completeness of the sem-usage.
// However, in production, we create sem without the need of closing the sem.
//
// We simply use ipcrm to delete the sem if necessary.
//
// Currently used only in test.
//
// XXX [2020-12-06] We don't do PasswdDestroy.
//
// Just let PasswdInit do the checking to avoid
// the duplication of sem.
func PasswdDestroy() error {
	if !IsTest {
		return ErrInvalidOp
	}

	if types.IS_ALL_GUEST {
		return nil
	}

	if Sem == nil {
		return ErrSemNotExists
	}

	err := Sem.Destroy(0)
	if err != nil {
		log.Errorf("cmbbs.PasswdDestroy: unable to close: e: %v", err)
		return err
	}

	Sem = nil

	log.Infof("cmbbs.PasswdDestroy: done")

	return nil
}
