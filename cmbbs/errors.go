package cmbbs

import "errors"

var (
	ErrSemAlreadyExists   = errors.New("sem already exists")
	ErrSemNotExists       = errors.New("sem not exists")
	ErrInvalidOp          = errors.New("invalid op")
	ErrInvalidPasswd2Size = errors.New("invalid passwd2 size")
)
