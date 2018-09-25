package impl

import (
	"github.com/dgraph-io/badger"
	"time"
	"github.com/gushitong/aryadb/io"
)

type badgerStorage struct{
	opt badger.Options
	DB *badger.DB
}

func (b badgerStorage) NewTransaction(update bool) io.Transaction {
	return &badgeTxn{b.DB.NewTransaction(update)}
}

func (b badgerStorage) View(fn func(txn io.Transaction) error) error {
	txn := b.NewTransaction(false)
	defer txn.Discard()
	return fn(txn)
}

func (b badgerStorage) Update(fn func(txn io.Transaction) error) error {
	txn := b.NewTransaction(true)
	defer txn.Discard()

	if err := fn(txn); err != nil {
		return err
	}
	return txn.Commit(nil)
}

type badgerIterator struct {
	*badger.Iterator
}

func (i badgerIterator) GetItem() io.Item {
	return i.Item()
}

type badgeTxn struct {
	Txn *badger.Txn
}

func (t badgeTxn) IncrBy(key []byte, v int64)(int64, error) {
	val, err := t.Get(key)
	if err != nil {
		return 0, err
	}

	var num int64
	if val == nil {
		num = 0
	}else {
		num, err = io.ParseInt64(val)
		if err != nil {
			return 0, err
		}
	}
	num += v
	return num, t.Set(key, io.Int642Byte(num))
}

func (t badgeTxn) IncrByFloat(key []byte, v float64)(float64, error) {
	val, err := t.Get(key)
	if err != nil {
		return 0, err
	}
	var num float64
	if val == nil {
		num = 0
	} else {
		num, err = io.ParseFloat64(val)
		if err != nil {
			return 0, err
		}
	}
	num += v
	return num, t.Set(key, io.Float642Byte(num))
}

func (t badgeTxn) Get(key []byte) ([]byte, error) {
	item, err := t.Txn.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}else {
			return nil, err
		}
	}
	return item.Value()
}

func (t badgeTxn) Set(key, val []byte) error {
	return t.Txn.Set(key, val)
}

func(t badgeTxn) SetWithTTL(key, val []byte, ttl time.Duration) error {
	return t.Txn.SetWithTTL(key, val, ttl)
}

func(t badgeTxn) Del(key []byte) error {
	return t.Txn.Delete(key)
}

func(t badgeTxn) NewIterator(ops io.IteratorOptions) io.Iterator {
	it := t.Txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: ops.PrefetchValues,
		PrefetchSize: ops.PrefetchSize,
		Reverse: ops.Reverse,
		AllVersions: ops.AllVersions,
	})
	return badgerIterator{it}
}

func(t badgeTxn) Commit(fn func(error)) error {
	return t.Txn.Commit(fn)
}

func(t badgeTxn) Discard() {
	t.Txn.Discard()
}

func NewBadgerStorage(dir, valueDir string) (*badgerStorage, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = valueDir
	bdg, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &badgerStorage{opt:opts, DB: bdg}, nil
}