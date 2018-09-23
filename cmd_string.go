package main

import (
	"github.com/gushitong/aryadb/io"
	"math/bits"
	"strconv"
)

func _append(db io.DB, conn Conn, req Request) {
	var n int
	err := db.Update(func(txn io.Transaction) error {
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

func bitcount(db io.DB, conn Conn, req Request) {
	var n int
	err := db.View(func(txn io.Transaction) error {
		val, err := txn.Get(req.Args[1])
		if err != nil {
			return err
		}

		if val == nil {
			n = 0
			return nil
		}

		for i := 0; i < len(val); i++ {
			n += bits.OnesCount8(val[i])
		}
		return nil
	})

	if err != nil {
		conn.WriteRawError(err)
		return
	}
	conn.WriteInt(n)
}

func decr(db io.DB, conn Conn, req Request) {
	var v int64
	err := db.Update(func(txn io.Transaction) error {
		n, err := txn.IncrBy(req.Args[1], -1)
		if err != nil {
			return err
		}
		v = n
		return nil
	})

	if err != nil {
		conn.WriteRawError(err)
		return
	}
	conn.WriteInt64(v)
}

func decrby(db io.DB, conn Conn, req Request) {
	var v int64
	err := db.Update(func(txn io.Transaction) error {
		num, err := strconv.ParseInt(string(req.Args[2]), 10, 64)
		if err != nil {
			return err
		}
		n, err := txn.IncrBy(req.Args[1], -1*num)
		if err != nil {
			return err
		}
		v = n
		return nil
	})
	if err != nil {
		conn.WriteRawError(err)
		return
	}
	conn.WriteInt64(v)
}

func ping(db io.DB, conn Conn, req Request) {
	conn.WriteString("PONG")
}

func get(db io.DB, conn Conn, req Request) {
	var v []byte

	err := db.View(func(txn io.Transaction) error {
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

func set(db io.DB, conn Conn, req Request) {
	err := db.Update(func(txn io.Transaction) error {
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
