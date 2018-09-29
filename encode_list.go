package main

import (
	"bytes"
	"encoding/binary"
	"github.com/gushitong/aryadb/stor"
	"github.com/gushitong/aryadb/ut"
	"io"
)

type ListOperation int8

const (
	Lpush ListOperation = iota
	Lpop
	Rpush
	Rpop
)

type ListEncoder struct {
	key []byte
}

func (l ListEncoder) Encode(seq int64) []byte {
	var flag byte
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, seq)
	if seq < 0 {
		flag = '-'
	} else {
		flag = '='
	}
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolList)},
		[]byte{uint8(len(l.key))},
		l.key,
		[]byte{flag},
		buf.Bytes(),
	)
}

func (l ListEncoder) EncodeMinSeqKey() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolListIndex)},
		[]byte{uint8(len(l.key))},
		l.key,
		[]byte{(byte)(SymbolMinSeq)},
	)
}

func (l ListEncoder) EncodeMaxSeqKey() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolListIndex)},
		[]byte{uint8(len(l.key))},
		l.key,
		[]byte{(byte)(SymbolMaxSeq)},
	)
}

func (l ListEncoder) QueueKey() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolQueue)},
		l.key,
	)
}

func (l ListEncoder) Prefix() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolList)},
		[]byte{uint8(len(l.key))},
		l.key,
	)
}

func (l ListEncoder) SeqPrefix() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolListIndex)},
		[]byte{uint8(len(l.key))},
		l.key,
	)
}

func (l ListEncoder) Meta(txn stor.Transaction) (queue int64, minSeq int64, err error) {
	queueVal, err := txn.Get(l.QueueKey())
	if err != nil {
		return
	}
	queue, err = ut.Bytes2Int64(queueVal)
	if err != nil && err != io.EOF {
		return
	}

	minSeqKey := l.EncodeMinSeqKey()
	seqVal, err := txn.Get(minSeqKey)
	if err != nil {
		return
	}
	minSeq, err = ut.Bytes2Int64(seqVal)
	if err == io.EOF {
		err = nil
	}
	return
}

func (l ListEncoder) UpdateCapacity(op ListOperation, txn stor.Transaction) (int64, error) {
	var queueLen, incr int64
	queueKey := l.QueueKey()
	if val, err := txn.Get(queueKey); err != nil {
		return 0, err
	} else if val == nil {
		queueLen = 0
	} else {
		if queueLen, err = ut.Bytes2Int64(val); err != nil {
			return 0, err
		}
	}

	switch op {
	case Lpush, Rpush:
		incr = 1
	case Lpop, Rpop:
		incr = -1
	}

	if queueLen+incr < 0 {
		return 0, ErrQueueEmpty
	}
	buf := ut.Int642Bytes(queueLen + incr)
	return queueLen + incr, txn.Set(queueKey, buf)
}

func (l ListEncoder) UpdateBoundary(op ListOperation, txn stor.Transaction) (int64, int64, error) {
	var err error
	var seqKey []byte
	var incr int64
	var newBoundary, oldBoundary int64

	switch op {
	case Lpush:
		incr = -1
		seqKey = l.EncodeMinSeqKey()
	case Lpop:
		incr = 1
		seqKey = l.EncodeMinSeqKey()
	case Rpush:
		incr = 1
		seqKey = l.EncodeMinSeqKey()
	case Rpop:
		incr = -1
		seqKey = l.EncodeMinSeqKey()
	}
	val, err := txn.Get(seqKey)
	if err != nil {
		return 0, 0, err
	} else if val != nil {
		oldBoundary, err = ut.Bytes2Int64(val)
		newBoundary = oldBoundary + incr
	}
	buf := ut.Int642Bytes(newBoundary)
	return oldBoundary, newBoundary, txn.Set(seqKey, buf)
}

func (l ListEncoder) UpdateMeta(op ListOperation, txn stor.Transaction) (int64, int64, int64, error) {
	i1, i2, err := l.UpdateBoundary(op, txn)
	if err != nil {
		return 0, 0, 0, err
	}
	q, err := l.UpdateCapacity(op, txn)
	if err != nil {
		return 0, 0, 0, err
	}
	return i1, i2, q, err
}

func (l ListEncoder) DecodeSeq(listKey []byte) (int64, error) {
	if listKey[0] != (byte)(SymbolList) {
		return 0, ErrCorruptedListKey
	}

	lenKey := int(listKey[1])
	if len(listKey) <= lenKey+3 {
		return 0, ErrCorruptedListKey
	}
	return ut.Bytes2Int64(listKey[lenKey+3:])
}

func NewListEncoder(key []byte) (*ListEncoder, error) {
	if len(key) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}
	return &ListEncoder{key: key}, nil
}
