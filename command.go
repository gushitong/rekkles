package main

import (
	"github.com/gushitong/aryadb/db"
)

func Append(db db.DB, conn Conn, req Request) {
	var n int
	err := db.ReadWrite(func(txn db.Transaction) error {
		val, err := txn.Get(req.Args[1])
		if err != nil {
			return err
		}
		val = append(val, req.Args[2]...)
		n = len(val)
		if err := txn.Set(req.Args[1], val); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		conn.WriteRawError(err)
		return
	}
	conn.WriteInt(n)
}

func ping(db db.DB, conn Conn, req Request) {
	conn.WriteString("PONG")
}

func get(db db.DB, conn Conn, req Request) {
	var v []byte

	err := db.Read(func(txn db.Transaction) error {
		val, err := txn.Get(req.Args[1])
		if err != nil {
			return err
		}
		v = val
		return nil
	})

	if err != nil {
		conn.WriteRawError(err)
		return
	}
	conn.WriteBulk(v)
}

func set(db db.DB, conn Conn, req Request) {
	err := db.ReadWrite(func(txn db.Transaction) error {
		if err := txn.Set(req.Args[1], req.Args[2]); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		conn.WriteRawError(err)
		return
	}
	conn.WriteString("OK")
}
