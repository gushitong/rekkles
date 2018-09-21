package main

import (
	"github.com/tidwall/redcon"
	"github.com/gushitong/aryadb/badger"
	"fmt"
)

type Handler func(s Storage, conn redcon.Conn, cmd redcon.Command)

type aryadb struct {

	storage Storage

	Handlers map[string]Handler
}

func (db *aryadb) Handle(conn redcon.Conn, cmd redcon.Command)  {
	handler, ok := db.Handlers[lower(cmd.Args[0])]
	if !ok {
		conn.WriteError(fmt.Sprintf("ERR command not supported: %s", string(cmd.Args[0])))
		return
	}
	handler(db.storage, conn, cmd)
}

func(db *aryadb) RegisterHandlers() {
	RegisterHandler(db,"get", get)
	RegisterHandler(db,"set", set)
	RegisterHandler(db, "ping", ping)
}


func NewAryaDB() *aryadb {
	storage, err := badger.NewBadgerStorage("/tmp/badger", "/tmp/badger")
	if err != nil{
		panic(err)
	}
	db := &aryadb{
		storage: storage, Handlers: make(map[string]Handler),
	}
	db.RegisterHandlers()
	return db
}
