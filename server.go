package main

import (
	"fmt"
	"github.com/gushitong/aryadb/impl"
	"github.com/gushitong/aryadb/io"
	"github.com/pkg/errors"
	"github.com/tidwall/redcon"
)

type Handler func(db io.DB, conn Conn, cmd Request)

type server struct {
	db       io.DB
	Auth     string
	Handlers map[string]Handler
}

func (s *server) RequirePass() bool {
	return s.Auth != ""
}

func (s *server) Authenticate(conn Conn, auth string) error {
	if s.Auth != auth {
		return errors.New("auth failed.")
	}
	ctx := conn.Context().(*Context)
	ctx.Authenticated = true
	return nil
}

func (s *server) Handle(redConn redcon.Conn, redCmd redcon.Command) {
	conn := Conn{redConn}
	req := Request{redCmd.Raw, redCmd.Args}

	if LowerString(req.Args[0]) == "auth" {
		if conn.Context() == nil {
			conn.SetContext(&Context{})
		}
		if err := s.Authenticate(conn, LowerString(req.Args[1])); err != nil {
			conn.SetContext(&Context{})
			conn.WriteError("ERR auth failed")
			return
		}
		conn.SetContext(&Context{Authenticated: true})
		conn.WriteString("OK")
		return
	} else if s.RequirePass() && conn.Authenticated() == false {

		conn.WriteString("ERR auth required")
		return
	} else if handler, ok := s.Handlers[LowerString(redCmd.Args[0])]; !ok {

		redConn.WriteError(fmt.Sprintf("ERR command not supported: %s", string(redCmd.Args[0])))
		return
	} else {

		handler(s.db, Conn{redConn}, Request{redCmd.Raw, redCmd.Args})
	}
}

func (s *server) RegisterHandlers() {
	RegisterHandler(s, "append", _append)
	RegisterHandler(s, "bitcount", bitcount)
	RegisterHandler(s, "decr", decr)
	RegisterHandler(s, "decrby", decrby)
	RegisterHandler(s, "get", get)
	RegisterHandler(s, "getbit", getbit)
	RegisterHandler(s, "set", set)
	//RegisterHandler(s, "ping", ping)
}

func NewAryadbServer() *server {
	storage, err := impl.NewBadgerStorage("/tmp/impl", "/tmp/impl")
	if err != nil {
		panic(err)
	}
	server := &server{
		db:       storage,
		Auth:     "",
		Handlers: make(map[string]Handler),
	}
	server.RegisterHandlers()
	return server
}
