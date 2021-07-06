package ptt

import "github.com/Ptt-official-app/go-pttbbs/ptttype"

const (
	//https://github.com/ptt/pttbbs/blob/master/mbbsd/acl.c
	BAKUMAN_OBJECT_TYPE_USER  = "u"
	BAKUMAN_OBJECT_TYPE_BOARD = "b"
)

var (
	BAKUMAN_REASON_LEN = ptttype.BTLEN
)
