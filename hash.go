package main

import (
	"github.com/gushitong/aryadb/stor"
	"github.com/gushitong/aryadb/ut"
)

func hdel(db stor.DB, conn aryConnection, cmd aryCommand) {
	var exists bool
	err := db.Update(func(txn stor.Transaction) error {
		key, err := cmd.HashKey()
		if err != nil {
			return err
		}

		if val, _ := txn.Get(key); val == nil {
			return nil
		} else {
			exists = true
		}
		return txn.Del(key)
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteBool(exists)
}

func hexists(db stor.DB, conn aryConnection, cmd aryCommand) {
	var exists bool
	err := db.View(func(txn stor.Transaction) error {
		key, err := cmd.HashKey()
		if err != nil {
			return err
		}
		if val, _ := txn.Get(key); val != nil {
			exists = true
		}
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteBool(exists)
}

func hget(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v []byte
	err := db.View(func(txn stor.Transaction) error {
		key, err := cmd.HashKey()
		if err != nil {
			return err
		}
		if val, err := txn.Get(key); err != nil {
			return err
		} else {
			v = val
			return nil
		}
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteBulk(v)
}

func hgetall(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := make([]string, 0)
	err := db.View(func(txn stor.Transaction) error {
		prefix, err := cmd.HashPrefix()
		if err != nil {
			return err
		}
		it := txn.NewIterator(stor.DefaultIteratorOptions)
		defer it.Close()
		encoder, err := NewHashEncoder(cmd.Args[0])
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.GetItem()
			hash, err := encoder.Decode(item.Key())
			if err != nil {
				return err
			}
			value, _ := item.Value()
			v = append(v, string(hash))
			v = append(v, string(value))
			//fmt.Println("[1] Len:", len(v))
			//fmt.Println("[1]", string(v[0]), string(v[1]))
		}
		//fmt.Println("[2] Len:", len(v))
		//fmt.Println("[2]", string(v[0]), string(v[1]))
		return nil
	})

	//fmt.Println("[3] Len:", len(v))
	//fmt.Println("[3]", string(v[0]), string(v[1]))
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteArray(len(v))
	for _, val := range v {
		conn.WriteString(val)
	}
}

func hincrby(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int64
	err := db.Update(func(txn stor.Transaction) error {
		key, err := cmd.HashKey()
		if err != nil {
			return err
		}
		n1, err := ut.ParseInt64(cmd.Args[2])
		if err != nil {
			return err
		}
		v, err = txn.IncrBy(key, n1)
		return err
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt64(v)
}

func hincrbyfloat(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v float64
	err := db.Update(func(txn stor.Transaction) error {
		key, err := cmd.HashKey()
		if err != nil {
			return err
		}
		n1, err := ut.ParseFloat64(cmd.Args[2])
		if err != nil {
			return err
		}
		v, err = txn.IncrByFloat(key, n1)
		return err
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteBulk(ut.Float642Byte(v))
}

func hkeys(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := make([]string, 0)
	err := db.View(func(txn stor.Transaction) error {
		encoder, _ := NewHashEncoder(cmd.Args[0])
		prefix := encoder.Prefix()
		it := txn.NewIterator(stor.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.GetItem()
			hash, err := encoder.Decode(item.Key())
			if err != nil {
				return err
			}
			v = append(v, string(hash))
		}
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteArray(len(v))
	for _, val := range v {
		conn.WriteString(val)
	}
}

func hlen(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int
	err := db.View(func(txn stor.Transaction) error {
		encoder, _ := NewHashEncoder(cmd.Args[0])
		prefix := encoder.Prefix()
		it := txn.NewIterator(stor.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.GetItem()
			_, err := encoder.Decode(item.Key())
			if err != nil {
				return err
			}
			v += 1
		}
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt(v)
}

func hmget(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := make([][]byte, 0)
	db.View(func(txn stor.Transaction) error {
		if len(cmd.Args) < 1 {
			return ErrWrongNumOfArguments
		}
		encoder, _ := NewHashEncoder(cmd.Args[0])
		for _, k := range cmd.Args[1:] {
			key := encoder.Encode(k)
			val, err := txn.Get(key)
			if err != nil {
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

func hmset(db stor.DB, conn aryConnection, cmd aryCommand) {
	err := db.Update(func(txn stor.Transaction) error {
		if len(cmd.Args) < 3 || len(cmd.Args)%2 != 1 {
			return ErrWrongNumOfArguments
		}
		encoder, err := NewHashEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		for i := 1; i < len(cmd.Args); i += 2 {
			key := encoder.Encode(cmd.Args[i])
			if err := txn.Set(key, cmd.Args[i+1]); err != nil {
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

func hscan(db stor.DB, conn aryConnection, cmd aryCommand) {
	conn.WriteErr(ErrCommandNotSupported)
}

func hset(db stor.DB, conn aryConnection, cmd aryCommand) {
	err := db.Update(func(txn stor.Transaction) error {
		key, err := cmd.HashKey()
		if err != nil {
			return err
		}
		return txn.Set(key, cmd.Args[2])
	})
	if err != nil {
		conn.WriteBool(false)
	} else {
		conn.WriteBool(true)
	}
}

func hsetnx(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v bool
	err := db.Update(func(txn stor.Transaction) error {
		key, err := cmd.HashKey()
		if err != nil {
			return err
		}
		if val, _ := txn.Get(key); val != nil {
			return nil
		}
		v = true
		return txn.Set(key, cmd.Args[2])
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteBool(v)
}

func hstrlen(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int
	db.View(func(txn stor.Transaction) error {
		key, _ := cmd.HashKey()
		if val, _ := txn.Get(key); val != nil {
			v = len(val)
		}
		return nil
	})
	conn.WriteInt(v)
}

func hvals(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := make([]string, 0)
	err := db.View(func(txn stor.Transaction) error {
		prefix, _ := cmd.HashPrefix()
		it := txn.NewIterator(stor.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.GetItem()
			val, _ := item.Value()
			v = append(v, string(val))
		}
		return nil
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteArray(len(v))
	for _, val := range v {
		conn.WriteString(val)
	}
}
