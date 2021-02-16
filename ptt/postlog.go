package ptt

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type PostLog struct {
	Author  ptttype.UserID_t
	Board   ptttype.BoardID_t
	Title   ptttype.Title_t
	TheDate types.Time4
	Number  int32
}

var (
	EMPTY_POSTLOG = PostLog{}
)

const (
	POSTLOG_SZ = unsafe.Sizeof(EMPTY_POSTLOG)
)
