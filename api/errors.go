package api

import "errors"

var (
	ErrInvalidHost       = errors.New("invalid host")
	ErrInvalidRemoteAddr = errors.New("invalid remote-addr")

	ErrInvalidParams = errors.New("invalid params")
	ErrInvalidPath   = errors.New("invalid path")
	ErrLoginFailed   = errors.New("login failed")
	ErrInvalidToken  = errors.New("invalid token")
)
