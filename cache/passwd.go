package cache

import (
	"encoding/binary"
	"os"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

//passwdUpdateMoney
//
//XXX should not call this directly.
//    call this from DeUMoney (SetUMoney).
func passwdUpdateMoney(uid ptttype.Uid, money int32) error {
	if uid < 1 || uid >= ptttype.MAX_USERS {
		return ErrInvalidUID
	}

	file, err := os.OpenFile(ptttype.FN_PASSWD, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	uidInCache := uid.ToUidInStore()
	const offsetMoney = unsafe.Offsetof(ptttype.USEREC_RAW.Money)
	file.Seek(int64(ptttype.USEREC_RAW_SZ*uintptr(uidInCache)+offsetMoney), 0)
	binary.Write(file, binary.LittleEndian, &money)

	return nil
}
