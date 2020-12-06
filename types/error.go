package types

import "errors"

var (
	ErrInvalidIni    = errors.New("invalid ini")
	ErrBytesTooLarge = errors.New("bytes too large")
	ErrNilReader     = errors.New("nil reader")
)
