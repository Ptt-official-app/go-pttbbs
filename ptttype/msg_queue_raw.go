package ptttype

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

type MsgQueueRaw struct {
	// Require updating SHM_VERSION if MSG_QUEUE_RAW_SZ is changed.
	Pid        types.Pid_t
	UserID     UserID_t
	LastCallIn [76]byte
	MsgMode    MsgMode
}

// Require updating SHM_VERSION if MSG_QUEUE_RAW_SZ is changed.
const MSG_QUEUE_RAW_SZ = unsafe.Sizeof(MsgQueueRaw{})
