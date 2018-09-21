package main

import (
	"github.com/tidwall/redcon"
	"github.com/gushitong/aryadb/impl"
	"fmt"
	"github.com/gushitong/aryadb/engine"
)

type Handler func(db engine.DB, conn Conn, cmd Request)

type db struct {
	db engine.DB

	Handlers map[string]Handler
}

func (db *db) Handle(conn redcon.Conn, cmd redcon.Command)  {
	handler, ok := db.Handlers[lower(cmd.Args[0])]
	if !ok {
		conn.WriteError(fmt.Sprintf("ERR command not supported: %s", string(cmd.Args[0])))
		return
	}
	handler(db.db, Conn{conn}, Request{cmd.Raw, cmd.Args})
}

func(db *db) RegisterHandlers() {
	RegisterHandler(db, "append", Append)
	RegisterHandler(db,"get", get)
	RegisterHandler(db,"set", set)
	RegisterHandler(db, "ping", ping)
}


func NewAryaDB() *db {
	storage, err := impl.NewBadgerStorage("/tmp/impl", "/tmp/impl")
	if err != nil{
		panic(err)
	}
	db := &db{
		db: storage, Handlers: make(map[string]Handler),
	}
	db.RegisterHandlers()
	return db
}
