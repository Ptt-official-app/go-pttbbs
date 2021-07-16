package types

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidIni    = errors.New("invalid ini")
	ErrBytesTooLarge = errors.New("bytes too large")
	ErrNilReader     = errors.New("nil reader")

	ErrNotImplemented = errors.New("not implemented")

	ErrInvalidTimeLocale = errors.New("invalid time locale")

	ErrInvalidFile = errors.New("invalid file")

	ErrInvalidSize = errors.New("invalid size")
)

func ErrRecover(err interface{}) error {
	str := fmt.Sprintf("(recover) %v", err)
	return errors.New(str)
}
