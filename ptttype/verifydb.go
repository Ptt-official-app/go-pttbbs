package ptttype

type VerifyDBStatus int8

const (
	VERIFYDB_OK    VerifyDBStatus = 0
	VERIFYDB_ERROR VerifyDBStatus = -1
)

type VerifyDBVMethod int32

const (
	VMETHOD_UNSET VerifyDBVMethod = iota
	VMETHOD_EMAIL
)

type VerifyDBMessage struct {
	RegMailDBReqHeader
	Message []byte
}
