package api

import "errors"

var (
	ErrInvalidParams = errors.New("invalid params")
	ErrInvalidPath   = errors.New("invalid path")
	ErrLoginFailed   = errors.New("login failed")
)
