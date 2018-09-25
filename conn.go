package main

import (
	"fmt"
	"github.com/tidwall/redcon"
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

func (c aryConnection) WriteBool(v bool)  {
	if v {
		c.WriteInt(1)
	} else {
		c.WriteInt(0)
	}
}