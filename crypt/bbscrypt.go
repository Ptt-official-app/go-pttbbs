package crypt

import (
	"errors"
)

var ErrInvalidCrypt = errors.New("invalid crypt")

// Fcrypt
//
// Params
//
//	key: the input-key (input-passwd) to be encrypted / checked
//	salt: the salt (expected-passwd-hash) in crypt(3)
//
// Return
//
//	theBytes: encrypted passwd, should be the same as salt if salt is the expected-passwd-hash.
//	err: err
func Fcrypt(key []byte, salt []byte) (theBytes []byte, err error) {
	passwdHash := [PASSLEN]byte{}
	cFcrypt(key, salt, &passwdHash)
	return passwdHash[:], nil
}
