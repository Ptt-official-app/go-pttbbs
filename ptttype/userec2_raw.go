package ptttype

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

const N_USEREC2_PAD_TAIL = DEFAULT_USEREC2_RAW_SZ -
	types.INT32_SZ - // Version
	PERM2_SZ - // UserLevel2
	types.TIME4_SZ

type Userec2Raw struct {
	Version uint32

	UserLevel2 PERM2
	UpdateTS   types.Time4

	PadTail [N_USEREC2_PAD_TAIL]byte
}

const DEFAULT_USEREC2_RAW_SZ = 128

var USEREC2_RAW = Userec2Raw{}

const USEREC2_RAW_SZ = unsafe.Sizeof(USEREC2_RAW)
