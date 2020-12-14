package bbs

import "errors"

var (
	ErrInvalidParams   = errors.New("invalid params")
	ErrInvalidBBoardID = errors.New("invalid bboardID")
)
