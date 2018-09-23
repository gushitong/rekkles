package main

import "github.com/pkg/errors"

var (
	ErrCommandNotSupported = errors.New("ERR command not supported")
	ErrKeyExists           = errors.New("ERR key exists")
	ErrBitValue            = errors.New("ERR bit is not an integer or out of range")
	ErrBitOffset           = errors.New("ERR bit offset is not an integer or out of range")
	ErrWrongNumOfArguments = errors.New("ERR wrong number of arguments.")
	ErrIntegerValue		= errors.New("ERR value is not an integer or out of range")
)
