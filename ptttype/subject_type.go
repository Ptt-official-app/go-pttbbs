package ptttype

//https://github.com/ptt/pttbbs/blob/master/include/common.h#L208
/* ----------------------------------------------------- */
/* 標題類形                                              */
/* ----------------------------------------------------- */

type SubjectType int

const (
	SUBJECT_NORMAL  SubjectType = 0
	SUBJECT_REPLY   SubjectType = 1
	SUBJECT_FORWARD SubjectType = 2
	SUBJECT_LOCKED  SubjectType = 3
	SUBJECT_DELETED SubjectType = 4
)
