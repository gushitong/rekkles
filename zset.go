package main

import (
	"github.com/gushitong/aryadb/stor"
	"github.com/gushitong/aryadb/ut"
	"bytes"
)

func zadd(db stor.DB, conn aryConnection, cmd aryCommand) {
	err := db.Update(func(txn stor.Transaction) error {
		if len(cmd.Args[0]) > MaxKeySize || len(cmd.Args[2]) > MaxKeySize {
			return ErrKeySizeExceeded
		}
		score, err := ut.ParseInt64(cmd.Args[1])
		if err != nil {
			return ErrIntegerValue
		}
		encoder, err := NewZsetEncoder(cmd.Args[0])
		if err != nil {
			return err
		}
		memberKey := encoder.EncodeMember(cmd.Args[2])
		if val, err := txn.Get(memberKey); err != nil {
			return err
		} else if val == nil {
			_, err = txn.IncrBy(encoder.QueueKey(), 1)
			if err != nil {
				return err
			}
		}
		err = txn.Set(memberKey, cmd.Args[1])
		if err != nil {
			return err
		}
		scoreKey := encoder.EncodeScore(score)
		return txn.Set(scoreKey, cmd.Args[2])
	})
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteInt64(1)
}

func zcard(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v int64
	err := db.View(func(txn stor.Transaction) error {
		encoder, _ := NewZsetEncoder(cmd.Args[0])
		val, err := txn.Get(encoder.QueueKey())
		if err != nil || val == nil {
			return err
		}
		v, err = ut.ParseInt64(val)
		return err
	})
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteInt64(v)
}

func zcount(db stor.DB, conn aryConnection, cmd aryCommand)  {
	var v int64
	err := db.View(func(txn stor.Transaction) error {
		encoder := ZsetEncoder{cmd.Args[0]}
		min, err := ut.ParseInt64(cmd.Args[1])
		if err != nil {
			return ErrIntegerValue
		}
		max, err := ut.ParseInt64(cmd.Args[2])
		if err != nil {
			return ErrIntegerValue
		}
		if max < min {
			return nil
		}
		minKey, maxKey := encoder.EncodeScore(min), encoder.EncodeScore(max)
		prefix := encoder.ScorePrefix()
		ops := stor.DefaultIteratorOptions
		ops.PrefetchValues = false
		it := txn.NewIterator(ops)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.GetItem().Key()
			if bytes.Compare(key, minKey) >= 0 || bytes.Compare(key, maxKey) <= 0 {
				v++
			} else if bytes.Compare(key, maxKey) > 0{
				break
			}
		}
		return nil
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteInt64(v)
}

func zincrby(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v []byte
	err := db.Update(func(txn stor.Transaction) error {
		encoder := &ZsetEncoder{cmd.Args[0]}
		incr, err := ut.ParseInt64(cmd.Args[1])
		if err != nil {
			return ErrIntegerValue
		}
		memberKey := encoder.EncodeMember(cmd.Args[2])
		val, err := txn.Get(memberKey)
		if err != nil {
			return err
		} else if val == nil {
			err := txn.Set(memberKey, cmd.Args[1])
			if err != nil {
				return err
			}
			scoreKey := encoder.EncodeScore(incr)
			err = txn.Set(scoreKey, cmd.Args[2])
			if err != nil {
				return err
			}
			v = ut.FormatInt64(incr)
		} else {
			score, err := ut.ParseInt64(val)
			if err != nil {
				return ErrCorruptedZsetScore
			}
			scoreKey := encoder.EncodeScore(score)
			err = txn.Del(scoreKey)
			if err != nil {
				return err
			}
			val = ut.FormatInt64(score+incr)
			err = txn.Set(memberKey, val)
			if err != nil {
				return err
			}
			scoreKey = encoder.EncodeScore(score+incr)
			err = txn.Set(scoreKey, cmd.Args[2])
			if err != nil {
				return err
			}
			v = ut.FormatInt64(incr+score)
		}
		return nil
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteBulk(v)
}

func zpopmax(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := make([]string, 0)
	err := db.View(func(txn stor.Transaction) error {
		encoder := &ZsetEncoder{cmd.Args[0]}
		prefix := encoder.ScorePrefix()
		ops := stor.DefaultIteratorOptions
		ops.PrefetchValues = false
		ops.Reverse = true
		it := txn.NewIterator(ops)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.GetItem()
			key := item.Key()
			value, err := item.Value()
			if err != nil {
				return err
			}
			score, err := encoder.DecodeScore(key)
			v = append(v, string(value))
			v = append(v, string(ut.FormatInt64(score)))
			break
		}
		_, err := txn.IncrBy(encoder.QueueKey(), -1)
		return err
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteArray(len(v))
	for _, val := range v {
		conn.WriteBulkString(val)
	}
}

func zpopmin(db stor.DB, conn aryConnection, cmd aryCommand) {
	v := make([]string, 0)
	err := db.View(func(txn stor.Transaction) error {
		encoder := &ZsetEncoder{cmd.Args[0]}
		prefix := encoder.ScorePrefix()
		ops := stor.DefaultIteratorOptions
		ops.PrefetchValues = false
		it := txn.NewIterator(ops)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.GetItem()
			key := item.Key()
			value, err := item.Value()
			if err != nil {
				return err
			}
			score, err := encoder.DecodeScore(key)
			v = append(v, string(value))
			v = append(v, string(ut.FormatInt64(score)))
			break
		}
		_, err := txn.IncrBy(encoder.QueueKey(), -1)
		return err
	})
	if err != nil {
		conn.WriteErr(err)
		return
	}
	conn.WriteArray(len(v))
	for _, val := range v {
		conn.WriteBulkString(val)
	}
}

func _zrange(db stor.DB, conn aryConnection, cmd aryCommand, reverse bool) {
	var v []string
	err := db.View(func(txn stor.Transaction) error {
		var withscores bool
		start, err := ut.ParseInt64(cmd.Args[1])
		if err != nil {
			return ErrIntegerValue
		}
		end, err := ut.ParseInt64(cmd.Args[2])
		if err != nil {
			return ErrIntegerValue
		}
		if len(cmd.Args) >= 4 && ut.LowerString(cmd.Args[3]) == "withscores" {
			withscores = true
		}
		encoder := &ZsetEncoder{cmd.Args[0]}
		queue, err := txn.Get(encoder.QueueKey())
		if err != nil {
			return err
		}
		queueLen, err := ut.ParseInt64(queue)
		if err != nil {
			return err
		}
		var i int64 = 0
		start = ut.FixBoundary(queueLen, start)
		end = ut.FixBoundary(queueLen, end)
		prefix := encoder.ScorePrefix()
		ops := stor.DefaultIteratorOptions
		ops.Reverse = reverse
		it := txn.NewIterator(stor.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.GetItem()
			scoreKey := item.Key()
			value, err := item.Value()
			if err != nil {
				return err
			}
			if i >= start && i <= end {
				v = append(v, string(value))
				if withscores {
					score, err := encoder.DecodeScore(scoreKey)
					if err != nil {
						return err
					}
					v = append(v, string(ut.FormatInt64(score)))
				}
			} else if i > end {
				break
			}
			i++
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

func zrange(db stor.DB, conn aryConnection, cmd aryCommand) {
	_zrange(db, conn, cmd, false)
}

func _zrangebyscore(db stor.DB, conn aryConnection, cmd aryCommand, reverse bool) {
	v := make([]string, 0)
	err := db.View(func(txn stor.Transaction) error {
		var withscores bool
		min, err := ut.ParseInt64(cmd.Args[1])
		if err != nil {
			return ErrIntegerValue
		}
		max, err := ut.ParseInt64(cmd.Args[2])
		if err != nil {
			return ErrIntegerValue
		}
		if max < min {
			return nil
		}
		encoder := ZsetEncoder{cmd.Args[0]}
		//minKey, maxKey := encoder.EncodeScore(min), encoder.EncodeScore(max)
		prefix := encoder.ScorePrefix()
		ops := stor.DefaultIteratorOptions
		ops.Reverse = reverse
		it := txn.NewIterator(ops)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.GetItem()
			scoreKey := item.Key()
			value, err := item.Value()
			if err != nil {
				return err
			}
			score, err := encoder.DecodeScore(scoreKey)
			if err != nil {
				return err
			}
			if min <= score && score <= max {
				v = append(v, string(value))
				if withscores {
					score, err := encoder.DecodeScore(scoreKey)
					if err != nil {
						return err
					}
					v = append(v, string(ut.FormatInt64(score)))
				}
			} else if score > max {
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

func zrangebyscore(db stor.DB, conn aryConnection, cmd aryCommand) {
	_zrangebyscore(db, conn, cmd, false)
}

func _zrank(db stor.DB, conn aryConnection, cmd aryCommand, reverse bool) {
	var v int64
	db.View(func(txn stor.Transaction) error {
		var i int64
		encoder := &ZsetEncoder{cmd.Args[0]}
		memberKey := encoder.EncodeMember(cmd.Args[1])
		val, err := txn.Get(memberKey)
		if err != nil {
			return err
		}
		if val == nil {
			return nil
		}
		score, err := ut.ParseInt64(val)
		if err != nil {
			return err
		}
		prefix := encoder.ScorePrefix()
		scoreKey := encoder.EncodeScore(score)
		ops := stor.DefaultIteratorOptions
		ops.Reverse = reverse
		ops.PrefetchValues = false
		it := txn.NewIterator(stor.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			i++
			if bytes.Compare(it.GetItem().Key(), scoreKey) == 0 {
				v = i
			}
		}
		return nil
	})
	if v == 0 {
		conn.WriteNull()
		return
	}
	conn.WriteInt64(v)
}

func zrank(db stor.DB, conn aryConnection, cmd aryCommand) {
	_zrank(db, conn, cmd, false)
}

func zrevrange(db stor.DB, conn aryConnection, cmd aryCommand) {
	_zrange(db, conn, cmd, true)
}

func zrevrangebyscore(db stor.DB, conn aryConnection, cmd aryCommand) {
	_zrangebyscore(db, conn, cmd, true)
}

func zrevrank(db stor.DB, conn aryConnection, cmd aryCommand) {
	_zrank(db, conn, cmd, true)
}

func zscore(db stor.DB, conn aryConnection, cmd aryCommand) {
	var v []byte
	err := db.View(func(txn stor.Transaction) error {
		encoder := &ZsetEncoder{cmd.Args[0]}
		memberKey := encoder.EncodeMember(cmd.Args[1])
		if val, err := txn.Get(memberKey); err != nil {
			return err
		} else {
			v = val
		}
		return nil
	})
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteBulk(v)
}