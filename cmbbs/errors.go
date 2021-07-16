package cmbbs

import "errors"

var (
	ErrSemAlreadyExists   = errors.New("sem already exists")
	ErrSemNotExists       = errors.New("sem not exists")
	ErrSemInvalid         = errors.New("sem invalid")
	ErrInvalidOp          = errors.New("invalid op")
	ErrInvalidPasswd2Size = errors.New("invalid passwd2 size")
)
