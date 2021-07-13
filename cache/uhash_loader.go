package cache

import (
	"io"
	"os"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"

	log "github.com/sirupsen/logrus"
)

//LoadUHash
//Load user-hash into SHM.
func LoadUHash() (err error) {
	if Shm == nil {
		return ErrShmNotInit
	}

	// line: 58
	number := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Number),
		unsafe.Sizeof(Shm.Raw.Number),
		unsafe.Pointer(&number),
	)

	loaded := int32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Loaded),
		unsafe.Sizeof(Shm.Raw.Loaded),
		unsafe.Pointer(&loaded),
	)

	//XXX in case it's not assumed zero, this becomes a race...
	if number == 0 && loaded == 0 {
		// line: 60
		err = fillUHash(false)
		if err != nil {
			return err
		}

		// line: 61
		// use golang style.
		todayIsZeroBytes := [ptttype.TODAYISSZ]byte{}
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.TodayIs),
			unsafe.Sizeof(Shm.Raw.TodayIs),
			unsafe.Pointer(&todayIsZeroBytes),
		)

		// line: 62
		loaded = 1
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.Loaded),
			unsafe.Sizeof(Shm.Raw.Loaded),
			unsafe.Pointer(&loaded),
		)
	} else {
		// line: 65
		err = fillUHash(true)
		if err != nil {
			return err
		}
	}

	return nil
}

var (
	uHashLoaderInvalidUserID = 0
)

func fillUHash(isOnfly bool) error {
	log.Infof("fillUHash: start: isOnfly: %v", isOnfly)
	InitFillUHash(isOnfly)

	file, err := os.Open(ptttype.FN_PASSWD)
	if err != nil {
		log.Errorf("fillUHash: unable to open passwd: file: %v e: %v", ptttype.FN_PASSWD, err)
		return err
	}
	defer file.Close()

	uidInCache := ptttype.UidInStore(0)

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

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Number),
		unsafe.Sizeof(Shm.Raw.Number),
		unsafe.Pointer(&uidInCache),
	)

	return nil
}

func userecRawAddToUHash(uidInCache ptttype.UidInStore, userecRaw *ptttype.UserecRaw, isOnfly bool) {
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

	shmUserID := ptttype.UserID_t{}
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Userid)+ptttype.USER_ID_SZ*uintptr(uidInCache),
		ptttype.USER_ID_SZ,
		unsafe.Pointer(&shmUserID),
	)

	offsetNextInHash := unsafe.Offsetof(Shm.Raw.NextInHash)

	if !isOnfly || types.Cstrcmp(userecRaw.UserID[:], shmUserID[:]) != 0 {
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.Userid)+ptttype.USER_ID_SZ*uintptr(uidInCache),
			ptttype.USER_ID_SZ,
			unsafe.Pointer(&userecRaw.UserID),
		)

		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.Money)+types.INT32_SZ*uintptr(uidInCache),
			types.INT32_SZ,
			unsafe.Pointer(&userecRaw.Money),
		)

		if ptttype.USE_COOLDOWN {
			zero := types.Time4(0)
			Shm.WriteAt(
				unsafe.Offsetof(Shm.Raw.CooldownTime)+types.TIME4_SZ*uintptr(uidInCache),
				types.TIME4_SZ,
				unsafe.Pointer(&zero),
			)
		}
	}

	p := h
	val := ptttype.UidInStore(0)
	offsetHashHead := unsafe.Offsetof(Shm.Raw.HashHead)
	//offsetNextInHash := unsafe.Offsetof(Shm.Raw.NextInHash)
	isFirst := true

	Shm.ReadAt(
		offsetHashHead+types.INT32_SZ*uintptr(p),
		types.INT32_SZ,
		unsafe.Pointer(&val),
	)

	l := 0
	for val >= 0 && val < ptttype.MAX_USERS {
		if isOnfly && val == uidInCache { // already in hash
			return
		}

		l++
		// go to next
		// 1. setting p as val
		// 2. get val from next_in_hash[p]
		p = uint32(val)
		Shm.ReadAt(
			offsetNextInHash+types.INT32_SZ*uintptr(p),
			types.INT32_SZ,
			unsafe.Pointer(&val),
		)

		isFirst = false
	}

	// set next in hash as n
	offset := offsetHashHead
	if !isFirst {
		offset = offsetNextInHash
	}
	val = uidInCache
	Shm.WriteAt(
		offset+types.INT32_SZ*uintptr(p),
		types.INT32_SZ,
		unsafe.Pointer(&val),
	)

	// set next in hash as -1
	p = uint32(val)
	val = -1
	Shm.WriteAt(
		offsetNextInHash+types.INT32_SZ*uintptr(p),
		types.INT32_SZ,
		unsafe.Pointer(&val),
	)
}

func InitFillUHash(isOnfly bool) {
	if !isOnfly {
		toFillHashHead := [1 << ptttype.HASH_BITS]ptttype.UidInStore{}
		for idx := range toFillHashHead {
			toFillHashHead[idx] = -1
		}
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.HashHead),
			unsafe.Sizeof(Shm.Raw.HashHead),
			unsafe.Pointer(&toFillHashHead),
		)
	} else {
		for idx := uint32(0); idx < (1 << ptttype.HASH_BITS); idx++ {
			checkHash(idx)
		}
	}
}

func checkHash(h uint32) {
	// p as delegate-pointer to the Shm.
	// in the beginning, p is the indicator of HashHead.
	// after 1st for-loop, p is in nextInHash.
	// val as the corresponding *p

	// line: 71
	p := h
	val := ptttype.UidInStore(0)
	pval := &val
	valptr := unsafe.Pointer(pval)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.HashHead)+types.INT32_SZ*uintptr(p),
		types.INT32_SZ,
		valptr,
	)

	// line: 72
	isFirst := true

	var offset uintptr
	offsetHashHead := unsafe.Offsetof(Shm.Raw.HashHead)
	offsetNextInHash := unsafe.Offsetof(Shm.Raw.NextInHash)

	userID := ptttype.UserID_t{}
	deep := 0
	for val != -1 {
		offset = offsetNextInHash
		if isFirst {
			offset = offsetHashHead
		}

		// check invalid pointer-val, set as -1  line: 74
		if val < -1 || val >= ptttype.MAX_USERS {
			log.Warnf("uhash_loader.checkHash: val invalid: h: %v p: %v val: %v isHead: %v", h, p, val, isFirst)
			*pval = -1
			Shm.WriteAt(
				offset+types.INT32_SZ*uintptr(p),
				types.INT32_SZ,
				valptr,
			)
			break
		}

		// get user-id: line: 75
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.Userid)+ptttype.USER_ID_SZ*uintptr(val),
			ptttype.USER_ID_SZ,
			unsafe.Pointer(&userID),
		)

		userIDHash := cmsys.StringHashWithHashBits(userID[:])

		// check hash as expected line: 76
		if userIDHash != h {
			// XXX
			// the result of the userID does not fit the h (broken?).
			// XXX uhash_loader is used only 1-time when starting the service.
			next := ptttype.UidInStore(0)

			// get next from *p (val)
			Shm.ReadAt(
				offsetNextInHash+types.INT32_SZ*uintptr(val),
				types.INT32_SZ,
				unsafe.Pointer(&next),
			)
			log.Warnf("userID hash is not in the corresponding idx (to remove) (%v): userID: %v userIDHash: %v h: %v next: %v", deep, types.CstrToString(userID[:]), userIDHash, h, next)
			// remove current by setting current as the next, hopefully the next user can fit the userIDHash.
			*pval = next
			Shm.WriteAt(
				offset+types.INT32_SZ*uintptr(p),
				types.INT32_SZ,
				unsafe.Pointer(&next),
			)
		} else {
			// 1. p as val (pointer in NextInHash)
			// 2. update val as NextInHash[p]
			p = uint32(val)
			Shm.ReadAt(
				offsetNextInHash+types.INT32_SZ*uintptr(p),
				types.INT32_SZ,
				unsafe.Pointer(&val),
			)
			isFirst = false

		}

		// line: 87
		deep++
		if deep == PRE_ALLOCATED_USERS+10 { //need to be larger than the pre-allocated users.
			//warn if it's too deep, we may need to consider enlarge the hash-table.
			log.Warnf("checkHash deep: %v h: %v p: %v val: %v isFirst: %v", deep, h, p, val, isFirst)
		}
	}
}
