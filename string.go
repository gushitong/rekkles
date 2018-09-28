package main

import (
	"github.com/gushitong/aryadb/stor"
	"github.com/gushitong/aryadb/ut"
	"math/bits"
	"strconv"
	"time"
)

func _append(db stor.DB, conn aryConnection, cmd aryCommand) {
	var n int
	err := db.Update(func(txn stor.Transaction) error {
		val, err := txn.Get(cmd.StringKey())
		if err != nil {
			return err
		}

		val = append(val, cmd.Args[1]...)
		n = len(val)
		if err := txn.Set(cmd.StringKey(), val); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt(n)
}

func bitcount(db stor.DB, conn aryConnection, cmd aryCommand) {
	var n int
	err := db.View(func(txn stor.Transaction) error {
		val, err := txn.Get(cmd.StringKey())
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
		conn.WriteErr(err)
		return
	}
	conn.WriteInt(n)
}

func decr(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int64
	err := db.Update(func(txn stor.Transaction) error {
		n, err := txn.IncrBy(cmd.StringKey(), -1)
		if err != nil {
			return err
		}
		v = n
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt64(v)
}

func decrby(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int64
	err := db.Update(func(txn stor.Transaction) error {
		num, err := strconv.ParseInt(string(cmd.Args[1]), 10, 64)
		if err != nil {
			return err
		}
		n, err := txn.IncrBy(cmd.StringKey(), -1*num)
		if err != nil {
			return err
		}
		v = n
		return nil
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt64(v)
}

func get(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v []byte
	err := db.View(func(txn stor.Transaction) error {
		val, err := txn.Get(cmd.StringKey())
		if err != nil {
			return err
		}
		v = val
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteBulk(v)
}

func getbit(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int
	err := db.View(func(txn stor.Transaction) error {
		val, err := txn.Get(cmd.StringKey())
		if err != nil {
			return err
		}
		i, err := strconv.Atoi(string(cmd.Args[1]))
		if err != nil {
			return err
		}
		if val == nil || len(val)*8 < i {
			v = 0
			return nil
		}
		v = int(val[i/8]>>uint(i%8)) & 1
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt(v)
}

func getrange(db stor.DB, conn aryConnection, cmd aryCommand) {
	var s string

	err := db.View(func(txn stor.Transaction) error {
		val, err := txn.Get(cmd.StringKey())
		if err != nil {
			return err
		}
		start, err := ut.ParseInt64(cmd.Args[1])
		if err != nil {
			return err
		}
		end, err := ut.ParseInt64(cmd.Args[2])
		if err != nil {
			return err
		}
		l := int64(len(val))
		o1 := ut.FixBoundary(l, start)
		o2 := ut.FixBoundary(l, end)
		if o2 == l-1 {
			s = string(val[o1:])
		} else {
			s = string(val[o1:o2]) + string(val[o2])
		}
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteString(s)
}

func getset(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v string
	err := db.Update(func(txn stor.Transaction) error {
		val, err := txn.Get(cmd.StringKey())
		if err != nil {
			return err
		}
		v = string(val)
		return txn.Set(cmd.StringKey(), cmd.Args[1])
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteString(v)
}

func incr(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int64
	err := db.Update(func(txn stor.Transaction) error {
		val, err := txn.IncrBy(cmd.StringKey(), 1)
		if err != nil {
			return err
		}
		v = val
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt64(v)
}

func incrby(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int64
	err := db.Update(func(txn stor.Transaction) error {
		n, err := strconv.ParseInt(string(cmd.Args[1]), 10, 64)
		if err != nil {
			return err
		}
		val, err := txn.IncrBy(cmd.StringKey(), n)
		if err != nil {
			return err
		}
		v = val
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt64(v)
}

func incrfloat(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v float64
	err := db.Update(func(txn stor.Transaction) error {
		n1, err := strconv.ParseFloat(string(cmd.Args[1]), 64)
		if err != nil {
			return err
		}
		val, err := txn.Get(cmd.StringKey())
		if err != nil {
			return err
		}
		n2, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return err
		}
		v = n1 + n2
		return txn.Set(cmd.Args[1], []byte(strconv.FormatFloat(v, 'f', -1, 64)))
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
}

func mget(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := make([][]byte, 0)
	db.View(func(txn stor.Transaction) error {
		for _, key := range cmd.Args {
			if val, err := txn.Get(NewStringEncoder(key).Encode()); err != nil {
				v = append(v, nil)
			} else {
				v = append(v, val)
			}
		}
		return nil
	})
	conn.WriteArray(len(v))
	for _, val := range v {
		conn.WriteBulk(val)
	}
}

func mset(db stor.DB, conn aryConnection, cmd aryCommand) {
	err := db.Update(func(txn stor.Transaction) error {
		if len(cmd.Args)%2 != 1 {
			return ErrWrongNumOfArguments
		}
		for i := 1; i < len(cmd.Args); i++ {
			if err := txn.Set(cmd.StringKey(), cmd.Args[i+1]); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteString("OK")
}

func msetnx(db stor.DB, conn aryConnection, cmd aryCommand) {
	err := db.Update(func(txn stor.Transaction) error {
		if len(cmd.Args)%2 != 1 {
			return ErrWrongNumOfArguments
		}
		for i := 1; i < len(cmd.Args); i += 2 {
			if val, err := txn.Get(cmd.StringKey()); err != nil {
				return err
			} else if val != nil {
				return ErrKeyExists
			}
		}
		for i := 1; i < len(cmd.Args); i++ {
			if err := txn.Set(cmd.StringKey(), cmd.Args[i+1]); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt(1)
}

func set(db stor.DB, conn aryConnection, cmd aryCommand) {
	err := db.Update(func(txn stor.Transaction) error {
		return txn.Set(cmd.StringKey(), cmd.Args[1])
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteString("OK")
}

func setbit(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := 0
	err := db.Update(func(txn stor.Transaction) error {
		b := 0
		if string(cmd.Args[2]) == "1" {
			b = 1
		} else if string(cmd.Args[2]) != "0" {
			return ErrBitValue
		}

		pos, err := strconv.Atoi(string(cmd.Args[1]))
		if err != nil {
			return ErrBitOffset
		}

		val, err := txn.Get(cmd.StringKey())
		if pos+1 > len(val)*8 {
			return nil
		}

		if b == 0 {
			val[pos/8] = ut.ClearBit(val[pos/8], uint(pos%8))
		} else {
			val[pos/8] = ut.SetBit(val[pos/8], uint(pos%8))
		}
		v = 1
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt(v)
}

func setex(db stor.DB, conn aryConnection, cmd aryCommand) {
	err := db.Update(func(txn stor.Transaction) error {
		val, err := strconv.Atoi(string(cmd.Args[1]))
		if err != nil {
			return ErrIntegerValue
		}
		return txn.SetWithTTL(cmd.StringKey(), cmd.Args[2], time.Duration(val)*time.Second)
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteString("OK")
}

func setnx(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v bool
	db.Update(func(txn stor.Transaction) error {
		if val, _ := txn.Get(cmd.StringKey()); val != nil {
			return nil
		} else {
			v = true
			return txn.Set(cmd.StringKey(), cmd.Args[1])
		}
	})
	conn.WriteBool(v)
}

func setrange(db stor.DB, conn aryConnection, cmd aryCommand) {
	conn.WriteErr(ErrCommandNotSupported)
}

func strlen(db stor.DB, conn aryConnection, cmd aryCommand) {
	var n int
	db.View(func(txn stor.Transaction) error {
		if val, _ := txn.Get(cmd.StringKey()); val == nil {
			n = 0
		} else {
			n = len(val)
		}
		return nil
	})
	conn.WriteInt(n)
}
