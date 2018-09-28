package main

import (
	"bytes"
	"encoding/binary"
	"github.com/gushitong/aryadb/stor"
	"github.com/gushitong/aryadb/ut"
)

func lindex(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v []byte
	err := db.View(func(txn stor.Transaction) error {
		index, err := ut.ParseInt64(cmd.Args[1])
		if err != nil {
			return err
		}
		encoder, err := NewListEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		queueLen, minSeqVal, err := encoder.Meta(txn)
		if err != nil {
			return err
		}
		if index >= queueLen {
			return nil
		}
		buf, _ := ut.Int642Bytes(minSeqVal + index)
		key := encoder.EncodeKey(buf)
		if val, err := txn.Get(key); err != nil {
			return err
		} else {
			v = val
			return nil
		}
	})
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteBulk(v)
}

func llen(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int64
	db.View(func(txn stor.Transaction) error {
		encoder, _ := NewListEncoder(cmd.Args[0])
		queueKey := encoder.QueueKey()
		val, err := txn.Get(queueKey)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(val)
		binary.Read(buf, binary.LittleEndian, &v)
		return nil
	})
	conn.WriteInt64(v)
}

func lpop(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v []byte
	err := db.Update(func(txn stor.Transaction) error {
		encoder, err := NewListEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		i1, _, q, err := encoder.UpdateMeta(Lpop, txn)
		if err != nil {
			return err
		}
		if q < 0 {
			return ErrQueueEmpty
		}
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, i1)
		key := encoder.EncodeKey(buf.Bytes())
		v, err = txn.Get(key)
		if err != nil {
			return err
		}
		return txn.Del(key)
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteBulk(v)
}

func lpush(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int64
	err := db.Update(func(txn stor.Transaction) error {
		encoder, err := NewListEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		_, i2, q, err := encoder.UpdateMeta(Lpush, txn)
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, &i2)
		key := encoder.EncodeKey(buf.Bytes())
		v = q
		return txn.Set(key, cmd.Args[1])
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt64(v)
}

func lpushx(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int64
	db.Update(func(txn stor.Transaction) error {
		qKey := cmd.QueueKey()
		if val, err := txn.Get(qKey); err != nil {
			return err
		} else {
			q, err := ut.Bytes2Int64(val)
			if err != nil || q == 0 {
				return nil
			}
		}
		encoder, err := NewListEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		i1, _, q, err := encoder.UpdateMeta(Lpush, txn)
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, i1)
		key := encoder.EncodeKey(buf.Bytes())
		v = q
		return txn.Set(key, cmd.Args[1])
	})
	conn.WriteInt64(v)
}

func lrange(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := make([]string, 0)
	err := db.View(func(txn stor.Transaction) error {
		encoder, _ := NewListEncoder(cmd.Args[0])
		queueLen, minSeqVal, err := encoder.Meta(txn)
		if err != nil {
			return err
		}
		if queueLen == 0 {
			return nil
		}
		i1, err := ut.ParseInt64(cmd.Args[1])
		if err != nil {
			return err
		}
		i2, err := ut.ParseInt64(cmd.Args[2])
		if err != nil {
			return err
		}
		i1 = ut.FixBoundary(queueLen, i1)
		i2 = ut.FixBoundary(queueLen, i2)

		for i := i1; i <= i2; i++ {
			buf, _ := ut.Int642Bytes(i + minSeqVal)
			key := encoder.EncodeKey(buf)
			val, _ := txn.Get(key)
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

func lset(db stor.DB, conn aryConnection, cmd aryCommand) {
	err := db.Update(func(txn stor.Transaction) error {
		encoder, err := NewListEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		queueLen, minSeq, err := encoder.Meta(txn)
		if err != nil {
			return err
		}
		index, err := ut.ParseInt64(cmd.Args[1])
		if err != nil {
			return err
		}
		if index >= queueLen {
			return ErrIndexOutOfRange
		}
		seq, _ := ut.Int642Bytes(minSeq + index)
		key := encoder.EncodeKey(seq)
		return txn.Set(key, cmd.Args[2])
	})

	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteString("OK")
}
