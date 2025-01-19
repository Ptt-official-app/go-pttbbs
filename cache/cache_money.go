package cache

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

// SetUMoney
//
// XXX uid-in-cache = uid - 1
func SetUMoney(uid ptttype.UID, money int32) (int32, error) {
	if types.IS_ALL_GUEST {
		return 0, nil
	}

	uidInCache := uid.ToUIDInStore()
	SHM.Shm.Money[uidInCache] = money

	err := passwdUpdateMoney(uid, money)
	if err != nil {
		return money, err
	}

	return MoneyOf(uid), nil
}

// DeUMoney
//
// Add money to uid. (money can be >= 0 or < 0)
// Get current money and set the money by adding to current-money.
func DeUMoney(uid ptttype.UID, money int32) (int32, error) {
	if types.IS_ALL_GUEST {
		return 0, nil
	}

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

func MoneyOf(uid ptttype.UID) (money int32) {
	if types.IS_ALL_GUEST {
		return 0
	}

	uidInCache := uid.ToUIDInStore()
	return SHM.Shm.Money[uidInCache]
}
