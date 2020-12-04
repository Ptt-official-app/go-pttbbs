package cache

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

//SetUMoney
//
//XXX uid-in-cache = uid - 1
func SetUMoney(uid ptttype.Uid, money int32) (int32, error) {
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Money)+types.INT32_SZ*uintptr(uid-1),
		types.INT32_SZ,
		unsafe.Pointer(&money),
	)
	err := passwdUpdateMoney(uid, money)
	if err != nil {
		return money, err
	}

	return MoneyOf(uid), nil

}

//DeUMoney
//
//Add money to uid. (money can be >= 0 or < 0)
//Get current money and set the money by adding to current-money.
func DeUMoney(uid ptttype.Uid, money int32) (int32, error) {
	if uid <= 0 || uid > ptttype.MAX_USERS {
		log.Errorf("DeUMoney: uid is invalid: uid: %v money: %v", uid, money)
		return -1, ErrInvalidUID
	}

	currentMoney := MoneyOf(uid)
	if money < 0 && currentMoney < -money {
		return SetUMoney(uid, 0)
	}

	return SetUMoney(uid, currentMoney+money)
}

func MoneyOf(uid ptttype.Uid) (money int32) {
	uidInCache := uid.ToUidInStore()
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Money)+types.INT32_SZ*uintptr(uidInCache),
		types.INT32_SZ,
		unsafe.Pointer(&money),
	)

	return money
}
