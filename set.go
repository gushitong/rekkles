package main

import (
	"bytes"
	"github.com/gushitong/rekkles/stor"
	"github.com/gushitong/rekkles/ut"
)

func sadd(db stor.DB, conn aryConnection, cmd aryCommand) {
	var added int
	err := db.Update(func(txn stor.Transaction) error {
		if len(cmd.Args) < 2 {
			return ErrWrongNumOfArguments
		}
		encoder, err := NewSetEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		for _, member := range cmd.Args[1:] {
			memberKey := encoder.EncodeMember(member)
			if val, _ := txn.Get(memberKey); bytes.Compare(val, SetMemberValue) == 0 {
				continue
			}
			if err := txn.Set(memberKey, SetMemberValue); err != nil {
				return err

			}
			if _, err := txn.IncrBy(encoder.QueueKey(), 1); err != nil {
				return err
			}
			added += 1
		}
		return nil
	})
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteInt(added)
}

func scard(db stor.DB, conn aryConnection, cmd aryCommand) {
	var card int64
	err := db.View(func(txn stor.Transaction) error {
		encoder, err := NewSetEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		val, err := txn.Get(encoder.QueueKey())
		if err != nil {
			return err
		}
		card, err = ut.ParseInt64(val)
		return err
	})
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteInt64(card)
}

func sismember(db stor.DB, conn aryConnection, cmd aryCommand) {
	var member bool
	db.View(func(txn stor.Transaction) error {
		encoder, err := NewSetEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		if val, _ := txn.Get(encoder.EncodeMember(cmd.Args[1])); bytes.Compare(val, SetMemberValue) == 0 {
			member = true
		}
		return nil
	})
	conn.WriteBool(member)
}

func smembers(db stor.DB, conn aryConnection, cmd aryCommand) {
	members := make([]string, 0)
	err := db.View(func(txn stor.Transaction) error {
		encoder, err := NewSetEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		prefix := encoder.Prefix()
		ops := stor.DefaultIteratorOptions
		ops.PrefetchValues = false
		it := txn.NewIterator(ops)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			member, err := encoder.DecodeMember(it.GetItem().Key())
			if err != nil {
				return err
			}
			members = append(members, string(member))
		}
		return nil
	})
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteArray(len(members))
	for _, val := range members {
		conn.WriteBulkString(val)
	}
}

func spop(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := make([]string, 0)
	err := db.Update(func(txn stor.Transaction) error {
		var n int64
		var err error
		if len(cmd.Args) >= 2 {
			n, err = ut.ParseInt64(cmd.Args[1])
			if err != nil {
				return err
			}
		} else {
			n = 1
		}
		encoder, err := NewSetEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		prefix := encoder.Prefix()
		ops := stor.DefaultIteratorOptions
		ops.PrefetchValues = false
		it := txn.NewIterator(ops)
		defer it.Close()
		var i int64 = 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.GetItem().Key()
			member, err := encoder.DecodeMember(key)
			if err != nil {
				return err
			}
			if err := txn.Del(key); err != nil {
				return err
			}
			if _, err := txn.IncrBy(encoder.QueueKey(), -1); err != nil {
				return err
			}
			v = append(v, string(member))
			i++
			if i >= n {
				break
			}
		}
		return nil
	})
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteArray(len(v))
	for _, val := range v {
		conn.WriteBulkString(val)
	}
}
