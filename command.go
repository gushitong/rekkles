package main

import (
	"github.com/gushitong/aryadb/engine"
)


func Append(db engine.DB, conn Conn, req Request) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	val, err := txn.Get(req.Args[1])
	val = append(val, req.Args[2]...)
	err = txn.Set(req.Args[1], val)
	err = txn.Commit(nil)
	if err != nil {
		conn.WriteRawError(err)
		return
	}

	conn.WriteInt(len(val))
}


func ping(db engine.DB, conn Conn, req Request) {
	conn.WriteString("PONG")
}


func get(db engine.DB, conn Conn, req Request) {
	txn := db.NewTransaction(false)
	defer txn.Discard()

	val, err := txn.Get(req.Args[1])
	if err != nil {
		conn.WriteRawError(err)
	}else {
		conn.WriteBulk(val)
	}
}


func set(db engine.DB, conn Conn, req Request) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	err := txn.Set(req.Args[1], req.Args[2])
	if err != nil {
		conn.WriteRawError(err)
		return
	}

	err = txn.Commit(nil)
	if err != nil {
		conn.WriteRawError(err)
		return
	}
	conn.WriteString("OK")
}