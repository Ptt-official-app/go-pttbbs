package cmbbs

import "errors"

var (
	ErrSemAlreadyExists = errors.New("sem already exists")
	ErrInvalidOp        = errors.New("invalid op")
)
