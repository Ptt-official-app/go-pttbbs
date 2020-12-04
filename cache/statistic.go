package cache

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func StatInc(stats uintptr) error {
	err := validateStats(stats)
	if err != nil {
		return err
	}

	Shm.IncUint32(unsafe.Offsetof(Shm.Raw.Statistic) + types.UINT32_SZ*stats)

	return nil
}

func CleanStat() {
	in := [ptttype.STAT_MAX]uint32{}
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Statistic),
		unsafe.Sizeof(Shm.Raw.Statistic),
		unsafe.Pointer(&in),
	)
}

func ReadStat(stats uintptr) (uint32, error) {
	err := validateStats(stats)
	if err != nil {
		return 0, err
	}

	out := uint32(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Statistic)+types.UINT32_SZ*stats,
		types.UINT32_SZ,
		unsafe.Pointer(&out),
	)

	return out, nil
}

func validateStats(stats uintptr) error {
	if Shm == nil {
		return ErrShmNotInit
	}

	if Shm.Raw.Version != SHM_VERSION {
		return ErrShmVersion
	}
	if stats < 0 || stats >= ptttype.STAT_MAX {
		return ErrStats
	}

	return nil
}
