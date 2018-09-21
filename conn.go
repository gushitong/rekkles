package main

import (
	"github.com/tidwall/redcon"
	"fmt"
)

type Conn struct {
	redcon.Conn
}

func (c Conn) WriteRawError(err error)  {
	c.WriteError(fmt.Sprintf("ERR %s", err))
}

type Request struct {
	Raw []byte
	Args [][]byte
}



