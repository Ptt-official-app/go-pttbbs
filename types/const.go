package types

import (
	"time"
)

const (
	MIN_TO_US     = 60000000   // minute to micro-second
	TS_TO_NANO_TS = 1000000000 // second to nano-second

	SYS_BUFFER_SIZE = 8192

	DEFAULT_FOLDER_CREATE_PERM = 0o755
	DEFAULT_FILE_CREATE_PERM   = 0o644
)

var TIMEZONE, _ = time.LoadLocation(TIME_LOCATION)
