package main

import (
	"fmt"
	"github.com/tidwall/redcon"
	"github.com/gushitong/aryadb/ut"
)

type Context struct {
	Authenticated bool
	ConnectionID  int64
}

type aryConnection struct {
	redcon.Conn
}

func (c aryConnection) Authenticated() bool {
	return c.Context() != nil && c.Context().(*Context).Authenticated
}

func (c aryConnection) WriteErr(err error) {
	c.WriteError(fmt.Sprint(err))
}

func (c aryConnection) WriteBool(v bool) {
	if v {
		c.WriteInt(1)
	} else {
		c.WriteInt(0)
	}
}

type aryCommand struct {
	Raw  []byte
	Args [][]byte
}

func (aryCmd aryCommand) StringKey() []byte {
	return NewStringEncoder(aryCmd.Args[0]).Encode()
}

func (aryCmd aryCommand) HashKey() ([]byte, error) {
	encoder, err := NewHashEncoder(aryCmd.Args[0])
	return encoder.Encode(aryCmd.Args[1]), err
}

func (aryCmd aryCommand) HashPrefix() ([]byte, error) {
	encoder, err := NewHashEncoder(aryCmd.Args[0])
	return encoder.Prefix(), err
}

func (aryCmd aryCommand) QueueKey() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolQueue)},
		aryCmd.Args[0],
	)
}