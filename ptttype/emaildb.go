package ptttype

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

type EmailDBOp uint8

const (
	_ EmailDBOp = iota
	REGMAILDB_REQ_COUNT
	REGMAILDB_REQ_SET
	REGCHECK_REQ_AMBIGUOUS
	VERIFYDB_MESSAGE
)

type RegMailDBReq struct {
	RegMailDBReqHeader
	UserID UserID_t
	Email  Email_t
}

const REG_MAILDB_REQ_SZ = unsafe.Sizeof(RegMailDBReq{})

type RegMailDBReqHeader struct {
	Cb        types.Size_t
	Operation EmailDBOp
}
