package ptttype

type PagerMode uint8

const (
	PAGER_OFF PagerMode = iota
	PAGER_ON
	PAGER_DISABLE
	PAGER_ANTIWB
	PAGER_FRIENDONLY
)
