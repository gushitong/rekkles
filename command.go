package main

import (
	"github.com/gushitong/aryadb/arya"
)

func Append(db arya.DB, conn Conn, req Request) {
	var n int
	err := db.ReadWrite(func(txn arya.Transaction) error {
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

func ping(db arya.DB, conn Conn, req Request) {
	conn.WriteString("PONG")
}

func get(db arya.DB, conn Conn, req Request) {
	var v []byte

	err := db.Read(func(txn arya.Transaction) error {
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

func set(db arya.DB, conn Conn, req Request) {
	err := db.ReadWrite(func(txn arya.Transaction) error {
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
