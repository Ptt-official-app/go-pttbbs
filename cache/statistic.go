package cache

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func StatInc(stats ptttype.Stat) error {
	if types.IS_ALL_GUEST {
		return nil
	}

	err := validateStats(stats)
	if err != nil {
		return err
	}

	SHM.Shm.Statistic[stats]++

	return nil
}

func CleanStat() {
	if types.IS_ALL_GUEST {
		return
	}

	SHM.Shm.Statistic = [ptttype.STAT_MAX]uint32{}
}

func ReadStat(stats ptttype.Stat) (uint32, error) {
	if types.IS_ALL_GUEST {
		return 0, nil
	}

	err := validateStats(stats)
	if err != nil {
		return 0, err
	}

	return SHM.Shm.Statistic[stats], nil
}

func validateStats(stats ptttype.Stat) error {
	if types.IS_ALL_GUEST {
		return nil
	}

	if SHM == nil {
		return ErrShmNotInit
	}

	if SHM.Raw.Version != SHM_VERSION {
		return ErrShmVersion
	}
	if stats < 0 || stats >= ptttype.STAT_MAX {
		return ErrStats
	}

	return nil
}
