package main

import "github.com/pkg/errors"

var (
	ErrKeyEmpty            = errors.New("ERR key empty")
	ErrQueueEmpty          = errors.New("ERR pop from empty queue")
	ErrCorruptedSetKey     = errors.New("ERR corrupted set key")
	ErrCorruptedHashKey    = errors.New("ERR corrupted hash key")
	ErrCorruptedListKey    = errors.New("ERR corrupted list key")
	ErrCorruptedZsetScore  = errors.New("ERR corrupt zset score")
	ErrCorruptedZsetMember = errors.New("ERR corrupt zset member")
	ErrKeySizeExceeded     = errors.New("ERR key size exceeded length: 255")
	ErrIndexOutOfRange     = errors.New("ERR index out of range")
	ErrCommandNotSupported = errors.New("ERR aryCommand not supported")
	ErrKeyExists           = errors.New("ERR key exists")
	ErrBitValue            = errors.New("ERR bit is not an integer or out of range")
	ErrBitOffset           = errors.New("ERR bit offset is not an integer or out of range")
	ErrWrongNumOfArguments = errors.New("ERR wrong number of arguments.")
	ErrIntegerValue        = errors.New("ERR value is not an integer or out of range")
)
