package ptttype

import "errors"

var (
	ErrUserIDAlreadyExists = errors.New("user id already exists")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidFilename     = errors.New("invalid filename")
)
