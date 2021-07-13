package cache

import (
	"encoding/binary"
	"os"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//passwdUpdateMoney
//
//XXX should not call this directly.
//    call this from DeUMoney (SetUMoney).
func passwdUpdateMoney(uid ptttype.Uid, money int32) (err error) {
	if uid < 1 || uid >= ptttype.MAX_USERS {
		return ErrInvalidUID
	}

	file, err := os.OpenFile(ptttype.FN_PASSWD, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	uidInCache := uid.ToUidInStore()
	const offsetMoney = unsafe.Offsetof(ptttype.USEREC_RAW.Money)
	_, err = file.Seek(int64(ptttype.USEREC_RAW_SZ*uintptr(uidInCache)+offsetMoney), 0)
	if err != nil {
		return err
	}
	return types.BinaryWrite(file, binary.LittleEndian, &money)
}
