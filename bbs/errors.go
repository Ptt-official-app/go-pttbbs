package bbs

import "errors"

var (
	ErrInvalidParams     = errors.New("invalid params")
	ErrInvalidBBoardID   = errors.New("invalid bboardID")
	ErrInvalidUUserID    = errors.New("invalid uuserID")
	ErrInvalidArticleID  = errors.New("invalid articleID")
	ErrInvalidPermission = errors.New("invalid permission")
)
