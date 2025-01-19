package cache

import "time"

// ////////
// !!! We should have only 1 SHM.
var SHM *Shm

const (
	// from https://github.com/ptt/pttbbs/blob/master/include/pttstruct.h
	// commit: 6bdd36898bde207683a441cdffe2981e95de5b20
	SHM_VERSION = 4842

	PRE_ALLOCATED_USERS = 1000

	BOARD_ZERO_TOTAL_RECHECK_TIME_US = 600 * 1000000
	BOARD_TOTAL_RECHECK_TIME_US      = 3600 * 1000000

	BOARD_ZERO_BOTTOM_RECHECK_TIME_US = 3600 * 1000000
	BOARD_BOTTOM_RECHECK_TIME_US      = 2 * 3600 * 1000000

	CRON_RELOAD_BCACHE_DURATION = 24 * time.Hour
)

var MAP *Map
