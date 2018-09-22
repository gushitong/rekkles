package main

import (
	"fmt"
	"github.com/gushitong/aryadb/arya"
	"github.com/gushitong/aryadb/impl"
	"github.com/pkg/errors"
	"github.com/tidwall/redcon"
	"log"
)

type Handler func(db arya.DB, conn Conn, cmd Request)

type server struct {
	db       arya.DB
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

	log.Printf("requirepass: %v. authenticated: %v\n", s.RequirePass(), conn.Authenticated())

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
	RegisterHandler(s, "append", Append)
	RegisterHandler(s, "get", get)
	RegisterHandler(s, "set", set)
	RegisterHandler(s, "ping", ping)
}

func NewAryaDB() *server {
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
