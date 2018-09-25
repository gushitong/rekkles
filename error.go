package main

import "github.com/pkg/errors"

var (
	ErrCorruptedHashKey    = errors.New("ERR corrupted hash key")
	ErrKeySizeExceeded     = errors.New("ERR key size exceeded length: 255")
	ErrCommandNotSupported = errors.New("ERR aryCommand not supported")
	ErrKeyExists           = errors.New("ERR key exists")
	ErrKeyType             = errors.New("ERR key type")
	ErrHashValueNotInt     = errors.New("ERR hash value is not an integer")
	ErrBitValue            = errors.New("ERR bit is not an integer or out of range")
	ErrBitOffset           = errors.New("ERR bit offset is not an integer or out of range")
	ErrWrongNumOfArguments = errors.New("ERR wrong number of arguments.")
	ErrIntegerValue        = errors.New("ERR value is not an integer or out of range")
)

