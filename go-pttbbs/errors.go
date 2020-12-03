package main

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidIni   = errors.New("invalid ini")
)
