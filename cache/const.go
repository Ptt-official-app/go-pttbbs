package cache

// ////////
// !!! We should have only 1 Shm.
var Shm *SHM

const (
	// from https://github.com/ptt/pttbbs/blob/master/include/pttstruct.h
	// commit: 6bdd36898bde207683a441cdffe2981e95de5b20
	SHM_VERSION = 4842

	PRE_ALLOCATED_USERS = 1000
)
