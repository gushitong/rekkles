package main

type Storage interface {
	Get(key []byte) ([]byte, error)

	Set(key, value []byte) error

	//SetWithTTL(key, value []byte, ttl time.Duration) error
}
