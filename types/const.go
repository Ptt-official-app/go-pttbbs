package types

import (
	"time"
)

const (
	TS_TO_NANO_TS = 1000000000 //second to nano-second

	SYS_BUFFER_SIZE = 8192

	DEFAULT_FOLDER_CREATE_PERM = 0755
)

var (
	TIMEZONE, _ = time.LoadLocation(TIME_LOCATION)
)
