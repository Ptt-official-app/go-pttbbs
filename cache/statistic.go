package cache

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func StatInc(stats ptttype.Stat) error {
	err := validateStats(stats)
	if err != nil {
		return err
	}

	Shm.Shm.Statistic[stats]++

	return nil
}

func CleanStat() {
	Shm.Shm.Statistic = [ptttype.STAT_MAX]uint32{}
}

func ReadStat(stats ptttype.Stat) (uint32, error) {
	err := validateStats(stats)
	if err != nil {
		return 0, err
	}

	return Shm.Shm.Statistic[stats], nil
}

func validateStats(stats ptttype.Stat) error {
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
