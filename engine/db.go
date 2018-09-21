package engine

import "time"


type DB interface {

	NewTransaction(update bool) Transaction

	View(fn func(txn Transaction) error) error

	Update(fn func(txn Transaction) error) error
}


type Transaction interface {

	Get(key []byte) ([]byte, error)

	Set(key, value []byte) error

	SetWithTTL(key, value []byte, ttl time.Duration) error

	Commit(callback func(error)) error

	Discard()
}

