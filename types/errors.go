package types

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")

	ErrInvalidTimeLocale = errors.New("invalid time locale")

	ErrInvalidFile = errors.New("invalid file")
)
