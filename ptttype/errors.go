package ptttype

import "errors"

var (
	ErrUserIDAlreadyExists  = errors.New("user id already exists")
	ErrBoardIDAlreadyExists = errors.New("board id already exists")
	ErrInvalidUserID        = errors.New("invalid user id")
	ErrInvalidBoardID       = errors.New("invalid board id")
	ErrInvalidFilename      = errors.New("invalid filename")
	ErrInvalidBid           = errors.New("invalid bid")
	ErrInvalidAid           = errors.New("invalid aid")
	ErrInvalidIdx           = errors.New("invalid idx")
	ErrInvalidType          = errors.New("invalid type")

	ErrInvalidAllowRejectEmailOp = errors.New("invalid allow-reject-email op")
)
