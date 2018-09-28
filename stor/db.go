package stor

import "time"

type IteratorOptions struct {
	PrefetchValues bool
	PrefetchSize   int
	Reverse        bool
	AllVersions    bool
}

var DefaultIteratorOptions = IteratorOptions{
	PrefetchValues: true,
	PrefetchSize:   100,
	Reverse:        false,
	AllVersions:    false,
}

type DB interface {
	NewTransaction(update bool) Transaction

	View(fn func(txn Transaction) error) error

	Update(fn func(txn Transaction) error) error
}

type Transaction interface {
	Get(key []byte) ([]byte, error)

	Set(key, value []byte) error

	SetWithTTL(key, value []byte, ttl time.Duration) error

	Del(key []byte) error

	IncrBy(key []byte, v int64) (int64, error)

	IncrByFloat(key []byte, v float64) (float64, error)

	NewIterator(ops IteratorOptions) Iterator

	Commit(callback func(error)) error

	Discard()
}

type Iterator interface {
	GetItem() Item

	Rewind()

	Seek(key []byte)

	Valid() bool

	ValidForPrefix(prefix []byte) bool

	Next()

	Close()
}

type Item interface {
	Key() []byte

	Value() ([]byte, error)
}
