package cache

import (
	"io"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"

	log "github.com/sirupsen/logrus"
)

// LoadUHash
// Load user-hash into SHM.
func LoadUHash() (err error) {
	if Shm == nil {
		return ErrShmNotInit
	}

	// line: 58
	number := Shm.Shm.Number
	loaded := Shm.Shm.Loaded

	// XXX in case it's not assumed zero, this becomes a race...
	if number == 0 && loaded == 0 {
		// line: 60
		err = fillUHash(false)
		if err != nil {
			return err
		}

		// line: 61
		// use golang style.
		todayIsZeroBytes := [ptttype.TODAYISSZ]byte{}
		Shm.Shm.TodayIs = todayIsZeroBytes

		// line: 62
		Shm.Shm.Loaded = 1
	} else {
		// line: 65
		err = fillUHash(true)
		if err != nil {
			return err
		}
	}

	return nil
}

var uHashLoaderInvalidUserID = 0

func fillUHash(isOnfly bool) error {
	log.Infof("fillUHash: start: isOnfly: %v", isOnfly)
	InitFillUHash(isOnfly)

	file, err := os.Open(ptttype.FN_PASSWD)
	if err != nil {
		log.Errorf("fillUHash: unable to open passwd: file: %v e: %v", ptttype.FN_PASSWD, err)
		return err
	}
	defer file.Close()

	uidInCache := ptttype.UIDInStore(0)

	uHashLoaderInvalidUserID = 0
	log.Infof("fillUHash: to for-loop: MAX_USERS: %v", ptttype.MAX_USERS)
	for ; ; uidInCache++ {
		userecRaw, eachErr := ptttype.NewUserecRawWithFile(file)
		if eachErr != nil {
			// io.EOF is reading correctly to the end the file.
			if eachErr == io.EOF {
				break
			}

			err = eachErr
			break
		}

		userecRawAddToUHash(uidInCache, userecRaw, isOnfly)
	}

	if err != nil {
		log.Errorf("fillUHash: unable to read passwd: file: %v e: %v", ptttype.FN_PASSWD, err)
		return err
	}

	log.Infof("fillUHash: to write usernum: %v", uidInCache)

	Shm.Shm.Number = int32(uidInCache)

	return nil
}

func userecRawAddToUHash(uidInCache ptttype.UIDInStore, userecRaw *ptttype.UserecRaw, isOnfly bool) {
	// uhash use userid="" to denote free slot for new register
	// However, such entries will have the same hash key.
	// So we skip most of invalid userid to prevent lots of hash collision.
	if !userecRaw.UserID.IsValid() {
		// dirty hack, preserve few slot for new register
		uHashLoaderInvalidUserID++
		if uHashLoaderInvalidUserID > PRE_ALLOCATED_USERS {
			return
		}
	}

	h := cmsys.StringHashWithHashBits(userecRaw.UserID[:])

	shmUserID := &Shm.Shm.Userid[uidInCache]
	if !isOnfly || types.Cstrcmp(userecRaw.UserID[:], shmUserID[:]) != 0 {
		Shm.Shm.Userid[uidInCache] = userecRaw.UserID
		Shm.Shm.Money[uidInCache] = userecRaw.Money
		if ptttype.USE_COOLDOWN {
			Shm.Shm.CooldownTime[uidInCache] = 0
		}
	}

	p := h
	val := Shm.Shm.HashHead[p]
	// offsetNextInHash := unsafe.Offsetof(Shm.Raw.NextInHash)
	isFirst := true

	l := 0
	for val >= 0 && val < ptttype.MAX_USERS {
		if isOnfly && val == uidInCache { // already in hash
			return
		}

		l++
		// go to next
		// 1. setting p as val
		// 2. get val from next_in_hash[p]
		p = cmsys.Fnv32_t(val)
		val = Shm.Shm.NextInHash[p]

		isFirst = false
	}

	// set next in hash as n

	if isFirst {
		Shm.Shm.HashHead[p] = uidInCache
	} else {
		Shm.Shm.NextInHash[p] = uidInCache
	}
	// set next in hash as -1
	Shm.Shm.NextInHash[uidInCache] = -1
}

func InitFillUHash(isOnfly bool) {
	if !isOnfly {
		toFillHashHead := [1 << ptttype.HASH_BITS]ptttype.UIDInStore{}
		for idx := range toFillHashHead {
			toFillHashHead[idx] = -1
		}
		Shm.Shm.HashHead = toFillHashHead
	} else {
		for idx := cmsys.Fnv32_t(0); idx < (1 << ptttype.HASH_BITS); idx++ {
			checkHash(idx)
		}
	}
}

func checkHash(h cmsys.Fnv32_t) {
	// p as delegate-pointer to the Shm.
	// in the beginning, p is the indicator of HashHead.
	// after 1st for-loop, p is in nextInHash.
	// val as the corresponding *p

	// line: 71
	p := h
	val := Shm.Shm.HashHead[p]

	// line: 72
	isFirst := true
	userID := (*ptttype.UserID_t)(nil)
	deep := 0
	for val != -1 {
		// check invalid pointer-val, set as -1  line: 74
		if val < -1 || val >= ptttype.MAX_USERS {
			log.Warnf("uhash_loader.checkHash: val invalid: h: %v p: %v val: %v isHead: %v", h, p, val, isFirst)
			if isFirst {
				Shm.Shm.HashHead[p] = -1
			} else {
				Shm.Shm.NextInHash[p] = -1
			}
			break
		}

		// get user-id: line: 75
		userID = &Shm.Shm.Userid[val]
		userIDHash := cmsys.StringHashWithHashBits(userID[:])

		// check hash as expected line: 76
		if userIDHash != h {
			// XXX
			// the result of the userID does not fit the h (broken?).
			// XXX uhash_loader is used only 1-time when starting the service.
			next := Shm.Shm.NextInHash[val]

			// get next from *p (val)
			log.Warnf("userID hash is not in the corresponding idx (to remove) (%v): userID: %v userIDHash: %v h: %v next: %v", deep, types.CstrToString(userID[:]), userIDHash, h, next)
			// remove current by setting current as the next, hopefully the next user can fit the userIDHash.
			val = next
			if isFirst {
				Shm.Shm.HashHead[p] = next
			} else {
				Shm.Shm.NextInHash[p] = next
			}
		} else {
			// 1. p as val (pointer in NextInHash)
			// 2. update val as NextInHash[p]
			p = cmsys.Fnv32_t(val)
			val = Shm.Shm.NextInHash[p]
			isFirst = false
		}

		// line: 87
		deep++
		if deep == PRE_ALLOCATED_USERS+10 { // need to be larger than the pre-allocated users.
			// warn if it's too deep, we may need to consider enlarge the hash-table.
			log.Warnf("checkHash deep: %v h: %v p: %v val: %v isFirst: %v", deep, h, p, val, isFirst)
		}
	}
}
