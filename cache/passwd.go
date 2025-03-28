package cache

import (
	"encoding/binary"
	"os"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// passwdUpdateMoney
//
// XXX should not call this directly.
// call this from DeUMoney (SetUMoney).
func passwdUpdateMoney(uid ptttype.UID, money int32) (err error) {
	if types.IS_ALL_GUEST {
		return nil
	}

	if uid < 1 || uid >= ptttype.MAX_USERS {
		return ErrInvalidUID
	}

	file, err := os.OpenFile(ptttype.FN_PASSWD, os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	uidInCache := uid.ToUIDInStore()
	const offsetMoney = unsafe.Offsetof(ptttype.USEREC_RAW.Money)
	_, err = file.Seek(int64(ptttype.USEREC_RAW_SZ*uintptr(uidInCache)+offsetMoney), 0)
	if err != nil {
		return err
	}
	return types.BinaryWrite(file, binary.LittleEndian, &money)
}
