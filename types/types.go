package types

import (
	"unsafe"
)

type Pid_t int32

type InAddr_t uint32

type Key_t int

type Size_t uint32

var (
	VERSION     string
	GIT_VERSION string
)

const INT32_SZ = unsafe.Sizeof(int32(0))

const UINT32_SZ = unsafe.Sizeof(uint32(0))

const PID_SZ = unsafe.Sizeof(Pid_t(0))

const UINT8_SZ = unsafe.Sizeof(uint8(0))
