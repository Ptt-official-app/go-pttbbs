package types

import (
	"errors"
	"fmt"
)

var (
	ErrNotImplemented = errors.New("not implemented")

	ErrInvalidTimeLocale = errors.New("invalid time locale")

	ErrInvalidFile = errors.New("invalid file")

	ErrInvalidSize = errors.New("invalid size")
)

func ErrRecover(err interface{}) error {
	str := fmt.Sprintf("(recover) %v", err)
	return errors.New(str)
}
