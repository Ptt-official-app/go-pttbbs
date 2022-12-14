package cache

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type shmBoardByName []ptttype.BidInStore

func (s shmBoardByName) Len() int {
	return len(s)
}

func (s shmBoardByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s shmBoardByName) Less(i, j int) (isLess bool) {
	iBCache := Shm.Shm.BSorted[ptttype.BSORT_BY_NAME][i]
	jBCache := Shm.Shm.BSorted[ptttype.BSORT_BY_NAME][j]
	isLess = types.Cstrcasecmp(Shm.Shm.BCache[iBCache].Brdname[:], Shm.Shm.BCache[jBCache].Brdname[:]) < 0
	return isLess
}

type shmBoardByClass []ptttype.BidInStore

func (s shmBoardByClass) Len() int {
	return len(s)
}

func (s shmBoardByClass) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s shmBoardByClass) Less(i, j int) bool {
	iBCache := Shm.Shm.BSorted[ptttype.BSORT_BY_CLASS][i]
	jBCache := Shm.Shm.BSorted[ptttype.BSORT_BY_CLASS][j]

	cmp := types.Cstrcmp(Shm.Shm.BCache[iBCache].Title[:4], Shm.Shm.BCache[jBCache].Title[:4])
	if cmp != 0 {
		return cmp < 0
	} else {
		return types.Cstrcasecmp(Shm.Shm.BCache[iBCache].Brdname[:], Shm.Shm.BCache[jBCache].Brdname[:]) < 0
	}
}
