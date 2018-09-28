package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSetSadd(t *testing.T) {
	key := "sadd"
	val1, val2, val3 := "1", "2", "3"

	client.Del(key)

	n, err := client.SAdd(key, val1).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), n)

	n, err = client.SAdd(key, val2, val3).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), n)

	n, err = client.SAdd(key, val1, val2, val3).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), n)
}

func TestSetScard(t *testing.T) {
	key := "scard"
	val1, val2, val3 := "1", "2", "3"

	client.Del(key)
	client.SAdd(key, val2, val1, val3, val1)
	n, err := client.SCard(key).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), n)
}

func TestSetSismember(t *testing.T) {
	key := "scard"
	val1, val2, val3 := "1", "2", "3"
	client.Del(key)
	client.SAdd(key, val2, val1)
	ok, err := client.SIsMember(key, val1).Result()
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = client.SIsMember(key, val3).Result()
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestSetMembers(t *testing.T) {
	key := "scard"
	val1, val2, val3 := "1", "2", "3"
	client.Del(key)
	client.SAdd(key, val2, val1, val3)
	val, err := client.SMembers(key).Result()
	assert.Nil(t, err)
	assert.Len(t, val, 3)
}
