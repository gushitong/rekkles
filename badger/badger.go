package badger

import "github.com/dgraph-io/badger"

type badgerStorage struct{

	opts badger.Options

	db *badger.DB
}

func(s badgerStorage)Get(key []byte) ([]byte, error) {

	var value []byte
	var err error

	s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.Value()
		return nil
	})
	return value, err
}


func(s badgerStorage)Set(key, value []byte) error {

	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
	return err
}


func NewBadgerStorage(dir, valueDir string) (*badgerStorage, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = valueDir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &badgerStorage{opts:opts, db:db}, nil
}