package ptttype

import "unsafe"

type PERM2 uint32

const (
	PERM2_INVALID  PERM2 = 0o00000000000
	PERM2_ID_EMAIL PERM2 = 0o00000000001 /* 藍勾勾 email */
)

const (
	PERM2_SZ = unsafe.Sizeof(PERM2(0))
)
