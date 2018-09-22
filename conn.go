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

func (c Conn) WriteRawError(err error) {
	c.WriteError(fmt.Sprintf("ERR %s", err))
}

type Request struct {
	Raw  []byte
	Args [][]byte
}
