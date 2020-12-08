package main

import "errors"

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidIni        = errors.New("invalid ini")
	ErrInvalidHost       = errors.New("invalid host")
	ErrInvalidRemoteAddr = errors.New("invalid remote-addr")
)
