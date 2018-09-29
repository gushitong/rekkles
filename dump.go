package main

import (
	"fmt"
	"github.com/gushitong/aryadb/stor"
	"github.com/gushitong/aryadb/ut"
	"io"
	"bytes"
)

func DumpList(key []byte, txn stor.Transaction) {
	var queueLen int64

	encoder, _ := NewListEncoder(key)

	val, err := txn.Get(encoder.QueueKey())
	if err != nil {
		panic(err)
	}
	queueLen, err = ut.Bytes2Int64(val)
	if err != nil && err != io.EOF {
		panic(err)
	}

	dumpStr := fmt.Sprintf("[List]: %s ", string(key))
	prefix := encoder.Prefix()
	it := txn.NewIterator(stor.DefaultIteratorOptions)
	defer it.Close()

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.GetItem()
		listKey := item.Key()
		value, err := item.Value()
		if err != nil {
			panic(err)
		}
		i, err := encoder.DecodeSeq(listKey)
		if err != nil {
			panic(err)
		}
		dumpStr += fmt.Sprintf("%d %s ", i, string(value))
	}

	dumpStr += fmt.Sprintf("(len %d)", queueLen)
	fmt.Println(dumpStr)
}

func DumpSet(key []byte, txn stor.Transaction) {
	encoder, _ := NewSetEncoder(key)
	it := txn.NewIterator(stor.DefaultIteratorOptions)
	defer it.Close()

	dumpStr := fmt.Sprintf("[Queue]: %s ", string(key))
	for it.Seek(encoder.Prefix()); it.ValidForPrefix(encoder.Prefix()); it.Next() {
		item := it.GetItem()
		member, err := encoder.DecodeMember(item.Key())
		if err != nil {
			panic(err)
		}
		if value, _ := item.Value(); bytes.Compare(value, SetMemberValue) != 0 {
			panic("Err ser member value")
		}
		dumpStr += fmt.Sprintf(" %s", string(member))
	}

	queueLen, err := txn.Get(encoder.QueueKey())
	if err != nil {
		panic(err)
	}
	dumpStr += fmt.Sprintf(" (len %s)", queueLen)
	fmt.Println(dumpStr)
}

func DumpZset(key []byte, txn stor.Transaction) {
	encoder, _ := NewZsetEncoder(key)
	it := txn.NewIterator(stor.DefaultIteratorOptions)
	defer it.Close()

	dumpStr := fmt.Sprintf("[Zset]: %s ", string(key))
	for it.Seek(encoder.MemberPrefix()); it.ValidForPrefix(encoder.MemberPrefix()); it.Next() {
		item := it.GetItem()
		member, err := encoder.DecodeMemberKey(item.Key())
		if err != nil {
			panic(err)
		}
		if value, err := item.Value(); err != nil {
			panic("Err ser member value")
		} else {
			dumpStr += fmt.Sprintf("%s:%s ", string(member), string(value))
		}
	}

	queueLen, err := txn.Get(encoder.QueueKey())
	if err != nil {
		panic(err)
	}
	dumpStr += fmt.Sprintf(" (len %s)", queueLen)
	fmt.Println(dumpStr)
}


func DumpZsetScore(key []byte, txn stor.Transaction) {
	encoder, _ := NewZsetEncoder(key)
	it := txn.NewIterator(stor.DefaultIteratorOptions)
	defer it.Close()

	dumpStr := fmt.Sprintf("[Z]: %s ", string(key))
	for it.Seek(encoder.ScorePrefix()); it.ValidForPrefix(encoder.ScorePrefix()); it.Next() {
		item := it.GetItem()
		score, err := encoder.DecodeScoreKey(item.Key())
		if err != nil {
			panic(err)
		}
		if value, err := item.Value(); err != nil {
			panic("Err zset member")
		} else {
			dumpStr += fmt.Sprintf("%d:%s ", score, string(value))
		}
	}

	queueLen, err := txn.Get(encoder.QueueKey())
	if err != nil {
		panic(err)
	}
	dumpStr += fmt.Sprintf(" (len %s)", queueLen)
	fmt.Println(dumpStr)
}