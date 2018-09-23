package main

import (
	"fmt"
	"github.com/tidwall/redcon"
)

type Context struct {
	Authenticated bool
	ConnectionID  int64
}

type Conn struct {
	redcon.Conn
}

func (c Conn) Authenticated() bool {
	return c.Context() != nil && c.Context().(*Context).Authenticated
}

func (c Conn) WriteErr(err error) {
	c.WriteError(fmt.Sprintf("ERR %s", err))
}

func (c Conn) WriteBool(v bool)  {
	if v {
		c.WriteInt(1)
	} else {
		c.WriteInt(0)
	}
}

type Request struct {
	Raw  []byte
	Args [][]byte
}
