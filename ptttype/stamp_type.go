package ptttype

type StampType uint8

//https://github.com/ptt/pttbbs/blob/master/common/bbs/fhdr_stamp.c#L26
const (
	STAMP_FILE StampType = 0
	STAMP_DIR  StampType = 1
	STAMP_LINK StampType = 2
)
