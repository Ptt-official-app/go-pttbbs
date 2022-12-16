package cache

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func GetUTMPNumber() (total int32) {
	return Shm.Shm.UTMPNumber
}

func SearchUListUserID(userID *ptttype.UserID_t) (utmpID ptttype.UtmpID, uInfo *ptttype.UserInfoRaw) {
	// start and end
	start := int32(0)

	end := Shm.Shm.UTMPNumber
	end--
	if end < 0 {
		return -1, nil
	}

	// current-sorted (for double-buffer)
	currentSorted := Shm.Shm.CurrSorted

	// search
	userIDInCache := (*ptttype.UserID_t)(nil)
	for i := (start + end) / 2; ; i = (start + end) / 2 {
		// get utmpID
		utmpID = Shm.Shm.Sorted[currentSorted][ptttype.SORT_BY_ID][i]

		// get user-id
		userIDInCache = &Shm.Shm.UInfo[utmpID].UserID

		// cmp
		j := types.Cstrcasecmp(userID[:], userIDInCache[:])

		if j == 0 {
			uInfo = &Shm.Shm.UInfo[utmpID]
			return utmpID, uInfo
		}

		// determine start / end
		if end == start {
			break
		} else if i == start {
			i = end //nolint
			start = end
		} else if j > 0 {
			start = i
		} else {
			end = i
		}
	}

	return -1, nil
}

func SearchUListPID(pid types.Pid_t) (utmpID ptttype.UtmpID, uInfo *ptttype.UserInfoRaw) {
	// start and end
	start := int32(0)

	end := Shm.Shm.UTMPNumber
	end--
	if end < 0 {
		return -1, nil
	}

	// current-sorted (for double-buffer)
	currentSorted := Shm.Shm.CurrSorted

	// search
	isDiff := types.Pid_t(0)
	uPid := types.Pid_t(0)
	for i := (start + end) / 2; ; i = (start + end) / 2 {
		// get utmpID
		utmpID = Shm.Shm.Sorted[currentSorted][ptttype.SORT_BY_PID][i]

		// get uPid
		uPid = Shm.Shm.UInfo[utmpID].Pid

		// do cmp()
		isDiff = pid - uPid

		if isDiff == 0 {
			uInfo = &Shm.Shm.UInfo[utmpID]
			return utmpID, uInfo
		}

		// determine start / end
		if end == start {
			break
		} else if i == start {
			i = end //nolint
			start = end
		} else if isDiff > 0 {
			start = i
		} else {
			end = i
		}
	}

	return -1, nil
}

// SetUtmpMode
//
// XXX skip utmp for now.
func SetUtmpMode(uid ptttype.UID, mode ptttype.UserOpMode) (err error) {
	pid := uid.ToPid()
	utmpID, _ := SearchUListPID(pid)
	if utmpID == -1 {
		return ErrNotFound
	}
	Shm.Shm.UInfo[utmpID].Mode = mode

	return nil
}
