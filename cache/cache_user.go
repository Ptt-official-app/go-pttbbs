package cache

import (
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"

	log "github.com/sirupsen/logrus"
)

// AddToUHash
func AddToUHash(uidInCache ptttype.UIDInStore, userID *ptttype.UserID_t) error {
	if types.IS_ALL_GUEST {
		return nil
	}

	h := cmsys.StringHashWithHashBits(userID[:])

	// line: 166
	SHM.Shm.Userid[uidInCache] = *userID

	// init vars
	p := h
	val := SHM.Shm.HashHead[p]

	// line: 168
	times := 0
	isNext := false
	for ; times < ptttype.MAX_USERS && val != -1; times++ {
		isNext = true
		p = cmsys.Fnv32_t(val)
		val = SHM.Shm.NextInHash[p]
	}
	if times >= ptttype.MAX_USERS {
		log.Errorf("Unable to add-to-uhash! uid-in-cache: %v userID: %v", uidInCache, string(userID[:]))
		return ErrAddToUHash
	}

	// set current ptr
	if !isNext {
		SHM.Shm.HashHead[p] = uidInCache
	} else {
		SHM.Shm.NextInHash[p] = uidInCache
	}

	// set next as -1
	SHM.Shm.NextInHash[uidInCache] = -1

	return nil
}

// RemoveFromUHash
func RemoveFromUHash(uidInCache ptttype.UIDInStore) error {
	if types.IS_ALL_GUEST {
		return nil
	}

	userID := &SHM.Shm.Userid[uidInCache]

	h := cmsys.StringHashWithHashBits(userID[:])

	// line: 191
	p := h
	val := SHM.Shm.HashHead[p]

	// line: 194
	times := 0
	isNext := false
	for ; times < ptttype.MAX_USERS && val != -1 && val != uidInCache; times++ {
		p = cmsys.Fnv32_t(val)
		isNext = true
		val = SHM.Shm.NextInHash[p]
	}
	if times >= ptttype.MAX_USERS {
		log.Errorf("Unable to remove-from-uhash! uid-in-cache: %v userID: %v", uidInCache, string(userID[:]))
		return ErrRemoveFromUHash
	}

	if val == uidInCache {
		nextNum := SHM.Shm.NextInHash[uidInCache]
		if !isNext {
			SHM.Shm.HashHead[p] = nextNum
		} else {
			SHM.Shm.NextInHash[p] = nextNum
		}
	}
	return nil
}

// SearchUser
// Params
//
//	userID: querying user-id.
//	isReturn: is return the user-id in the shm.
//
// Return
//
//	uid:
//	err:
func SearchUserRaw(userID *ptttype.UserID_t, rightID *ptttype.UserID_t) (uid ptttype.UID, err error) {
	if types.IS_ALL_GUEST {
		if rightID != nil {
			copy(rightID[:], ptttype.GUEST.UserID[:])
		}
		return ptttype.GUEST_UID, nil
	}

	if userID[0] == 0 {
		return 0, nil
	}
	return DoSearchUserRaw(userID, rightID)
}

func DoSearchUserRaw(userID *ptttype.UserID_t, rightID *ptttype.UserID_t) (ptttype.UID, error) {
	// XXX we should have 0 as non-exists.
	//     currently the reason why it's ok is because the probability of collision on 0 is low.
	if types.IS_ALL_GUEST {
		if rightID != nil {
			copy(rightID[:], ptttype.GUEST.UserID[:])
		}
		return ptttype.GUEST_UID, nil
	}

	_ = StatInc(ptttype.STAT_SEARCHUSER)
	h := cmsys.StringHashWithHashBits(userID[:])

	// p = SHM->hash_head[h]  //line: 219
	p := SHM.Shm.HashHead[h]

	shmUserID := (*ptttype.UserID_t)(nil)
	for times := 0; times < ptttype.MAX_USERS && p != -1 && p < ptttype.MAX_USERS; times++ {
		// if (strcasecmp(SHM->userid[p], userid) == 0)  //line: 222
		shmUserID = &SHM.Shm.Userid[p]
		if types.Cstrcasecmp(userID[:], shmUserID[:]) == 0 {
			if userID[0] != 0 && rightID != nil {
				copy(rightID[:], shmUserID[:])
			}
			return p.ToUID(), nil
		}
		p = SHM.Shm.NextInHash[p]
	}

	return 0, nil
}

// GetUserID
//
// XXX uid = uid-in-cache + 1
func GetUserID(uid ptttype.UID) (*ptttype.UserID_t, error) {
	if types.IS_ALL_GUEST {
		return &ptttype.GUEST.UserID, nil
	}

	uidInCache := uid.ToUIDInStore()
	if uidInCache < 0 || uidInCache >= ptttype.MAX_USERS {
		return nil, ErrInvalidUID
	}

	return &SHM.Shm.Userid[uidInCache], nil
}

// SetUserID
//
// XXX uid = uid-in-cache + 1
func SetUserID(uid ptttype.UID, userID *ptttype.UserID_t) (err error) {
	if types.IS_ALL_GUEST {
		return nil
	}

	if uid <= 0 || uid > ptttype.MAX_USERS {
		return ErrInvalidUID
	}

	uidInCache := uid.ToUIDInStore()
	errRemove := RemoveFromUHash(uidInCache)
	errAdd := AddToUHash(uidInCache, userID)
	if errRemove != nil {
		return errRemove
	}
	if errAdd != nil {
		return errAdd
	}

	return nil
}

// CooldownTimeOf
// https://github.com/ptt/pttbbs/blob/master/include/cmbbs.h#L97
func CooldownTimeOf(uid ptttype.UID) (cooldowntime types.Time4) {
	if types.IS_ALL_GUEST {
		return 0
	}

	uidInCache := uid.ToUIDInStore()

	// types.Time4 is int32, not uint32
	// we use 0x7FFFFFF0 instead of 0xFFFFFFF0
	return SHM.Shm.CooldownTime[uidInCache] & 0x7FFFFFF0
}

func AddCooldownTime(uid ptttype.UID, minutes int) (err error) {
	if types.IS_ALL_GUEST {
		return nil
	}

	cooldowntime := CooldownTimeOf(uid)
	base := types.NowTS()
	if base < cooldowntime {
		base = cooldowntime
	}

	base += types.Time4(minutes) * 60
	base &= 0x7FFFFFF0

	uidInCache := uid.ToUIDInStore()
	SHM.Shm.CooldownTime[uidInCache] = base

	return nil
}

func SetCooldownTime(uid ptttype.UID, cooldownTime types.Time4) (err error) {
	if types.IS_ALL_GUEST {
		return nil
	}

	uidInCache := uid.ToUIDInStore()
	SHM.Shm.CooldownTime[uidInCache] = cooldownTime
	return nil
}

// PosttimesOf
// https://github.com/ptt/pttbbs/blob/master/include/cmbbs.h#L98
func PosttimesOf(uid ptttype.UID) (posttimes types.Time4) {
	if types.IS_ALL_GUEST {
		return 0
	}

	uidInCache := uid.ToUIDInStore()

	return SHM.Shm.CooldownTime[uidInCache] & 0xF
}

func AddPosttimes(uid ptttype.UID, times int) (err error) {
	if types.IS_ALL_GUEST {
		return nil
	}

	posttimes := PosttimesOf(uid)
	newPosttimes := posttimes + types.Time4(times)

	uidInCache := uid.ToUIDInStore()
	if newPosttimes < 0xf {
		SHM.Shm.CooldownTime[uidInCache] += types.Time4(times)
	} else {
		SHM.Shm.CooldownTime[uidInCache] |= 0xf
	}

	return nil
}
