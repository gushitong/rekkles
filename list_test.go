package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListLindex(t *testing.T) {
	key := "lindex"
	val := "val"
	val2 := "val2"

	client.Del(key)
	client.LPush(key, val)
	v, err := client.LIndex(key, 0).Result()
	assert.Nil(t, err)
	assert.Equal(t, v, val)

	client.LPush(key, val2)
	v, err = client.LIndex(key, 0).Result()
	assert.Nil(t, err)
	assert.Equal(t, v, val2)
}

func TestListLlen(t *testing.T) {
	key := "llen"
	val := "val"

	client.Del(key)

	client.LPush(key, val)
	n, err := client.LLen(key).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), n)

	client.LPush(key, val)
	n, err = client.LLen(key).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), n)

	client.LPop(key)
	client.LPop(key)
	n, err = client.LLen(key).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), n)
}

func TestListLpush(t *testing.T) {
	key := "lpush"
	val1 := "val1"

	client.Del(key)

	n, err := client.LPush(key, val1).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), n)
	v, err := client.LIndex(key, 0).Result()
	assert.Nil(t, err)
	assert.Equal(t, val1, v)

	n, err = client.LPush(key, val1).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), n)
	v, err = client.LIndex(key, 1).Result()
	assert.Nil(t, err)
	assert.Equal(t, val1, v)

	n, err = client.LPush(key, val1).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), n)
	v, err = client.LIndex(key, 2).Result()
	assert.Nil(t, err)
	assert.Equal(t, val1, v)
}

func TestListLpushx(t *testing.T) {

	key := "lpushx"
	val := "val"

	client.Del(key)
	n, err := client.LPushX(key, val).Result()
	assert.Nil(t, err)
	assert.Zero(t, n)

	client.LPush(key, val)
	n, err = client.LPushX(key, val).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), n)
}

func TestListLrange(t *testing.T) {
	key := "lrange"
	val := "val"

	client.Del(key)
	v, err := client.LRange(key, 0, -1).Result()
	assert.Nil(t, err)
	assert.Len(t, v, 0)

	client.LPush(key, val)
	client.LPush(key, val)

	v, err = client.LRange(key, 0, -1).Result()
	assert.Nil(t, err)
	assert.Len(t, v, 2)
}

func TestListLset(t *testing.T) {
	key := "lset"
	val := "val"
	val2 := "val2"

	client.Del(key)
	client.LPush(key, val)
	client.LPush(key, val)
	ok, err := client.LSet(key, 1, val2).Result()
	assert.Nil(t, err)
	assert.Equal(t, "OK", ok)
	v, err := client.LIndex(key, 1).Result()
	assert.Nil(t, err)
	assert.Equal(t, v, val2)
}
