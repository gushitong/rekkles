package main

import (
	"github.com/gushitong/aryadb/stor"
	"github.com/gushitong/aryadb/ut"
)

func del(db stor.DB, conn aryConnection, cmd aryCommand) {
	if len(cmd.Args) == 0 {
		conn.WriteError(ErrWrongNumOfArguments.Error())
		return
	}
	err := db.Update(func(txn stor.Transaction) error {
		for _, key := range cmd.Args {
			// del string
			stringKey := NewStringEncoder(key).Encode()
			if err := txn.Del(stringKey); err != nil {
				return err
			}
			// del hash
			hashPrefix := ut.ConcatBytearray(
				[]byte{(byte)(SymbolHash)},
				[]byte{uint8(len(key))},
				key,
			)
			if err := delprefix(db, hashPrefix); err != nil {
				return err
			}
			// del set
			setPrefix := ut.ConcatBytearray(
				[]byte{(byte)(SymbolSet)},
				[]byte{uint8(len(key))},
				key,
			)
			if err := delprefix(db, setPrefix); err != nil {
				return err
			}
			// del zset
			zsetPrefix := ut.ConcatBytearray(
				[]byte{(byte)(SymbolZset)},
				[]byte{uint8(len(key))},
				key,
			)
			if err := delprefix(db, zsetPrefix); err != nil {
				return err
			}
			// del list
			listPrefix := ut.ConcatBytearray(
				[]byte{(byte)(SymbolList)},
				[]byte{uint8(len(key))},
				key,
			)
			if err := delprefix(db, listPrefix); err != nil {
				return err
			}
			// del list prefix
			listSeqPrefix := ut.ConcatBytearray(
				[]byte{(byte)(SymbolListIndex)},
				[]byte{uint8(len(key))},
				key,
			)
			if err := delprefix(db, listSeqPrefix); err != nil {
				return err
			}
			// del queue
			queueKey := ut.ConcatBytearray(
				[]byte{(byte)(SymbolQueue)},
				key,
			)
			if err := txn.Del(queueKey); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteString("OK")
}

func delprefix(db stor.DB, prefix []byte) error {
	return db.Update(func(txn stor.Transaction) error {
		ops := stor.DefaultIteratorOptions
		ops.PrefetchValues = false
		it := txn.NewIterator(ops)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.GetItem().Key()
			err := txn.Del(key)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func ping(db stor.DB, conn aryConnection, cmd aryCommand) {
	conn.WriteString("PONG")
}
